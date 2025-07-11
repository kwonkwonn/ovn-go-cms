package util

import "github.com/google/uuid"


func UUIDGenerator()(uuid.UUID,error){
	u,err:= uuid.NewRandom()
	if err!=nil{
		return  uuid.Nil,err
	}

	return u,nil
}