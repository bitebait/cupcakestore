package models

type Response struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Object  interface{} `json:"object"`
	Profile *Profile    `json:"profile"`
}
