package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	proto "github.com/contetto/user-service/proto"
	"github.com/gorilla/mux"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
)

var userClient proto.UsersClient

// swagger:route GET /user/{id} userAccountsCrudApi UserAPIFetch
//
// This will fetch an user.
//
// Will return user info in response
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Security:
//       api_key:
//
//     Responses:
//       200: user
//       500: internalError
//       400: validationError
//       401: unauthorisedError
func get(w http.ResponseWriter, r *http.Request) {
	// LH: @todo check authentification

	vars := mux.Vars(r)
	userID := vars["id"]

	isValid := bson.IsObjectIdHex(userID)
	if !isValid {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{"error": "Input validation failed. Invalid User ID provided"}`)
		return
	}

	// Call user-service to request the user object
	resp, err := userClient.Get(context.TODO(), &proto.GetReq{ID: userID})
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{"error": "Internal Error"}`)
		return
	}

	jsonRsp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{"error": "Error preparing result!"}`)
		return
	}

	w.WriteHeader(200)
	fmt.Fprint(w, string(jsonRsp))
}

func main() {
	// Create service
	service := micro.NewService(
		micro.Name("api-service"),
	)

	service.Init()

	// setup client to user-service
	userClient = proto.NewUsersClient("user-service", service.Client())

	r := mux.NewRouter()

	r.HandleFunc("/users/{id}", get).Methods("GET")
	http.Handle("/", r)
}
