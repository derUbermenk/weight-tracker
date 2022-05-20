import React, {Component, useEffect, useState} from "react";
import {useParams} from "react-router-dom";
import {RadioAttributes, TextAttributes, NumberAttributes, Button} from './UserComponents'

async function getUser(userId, userSetter) {
  const requestUrl = `http://localhost:8080/v1/api/user/${userId}`
  const requestOptions = {
    mode: 'cors'
  }

  const request = new Request(requestUrl, requestOptions)

  const response = await fetch(request);
  const user = await response.json();

  userSetter(user)
}

function User() {
  const {userId} = useParams(); 
  const [user, setUser] = useState({});

  useEffect(() => {
    getUser(userId, (user)=>setUser(user)) },
    []
  );
  
  const [isNotEdit, setisNotEdit] = useState(true);

  const handleUserChange = (attribute, value) => {
    var updated_user = user
    updated_user[attribute] = value
    setUser({...updated_user})
  }

  const handleEdit = () => {
    setisNotEdit(false) // make user editable
  }

  const handleSave = () => {
    alert(`Save Button clicked, saving ${JSON.stringify(user)}`)
    // saveUser()
    setisNotEdit(true) // make user non editable
  }

  return (
    <div>
      <h1>Hello {`${user.name}`}!</h1>
      <br></br>
      <div>
        {JSON.stringify(user)}
      </div>
      <br></br>
      { isNotEdit ? 
        <Button name={'Edit'} onclick={handleEdit}/> :
        <Button name={'Save'} onclick={handleSave}/> } 

      <form>
        <TextAttributes object={user} attribute="name" notEditable={isNotEdit} onChange={handleUserChange}/>
        <TextAttributes object={user} attribute="email" notEditable={isNotEdit} onChange={handleUserChange}/>
        <NumberAttributes object={user} attribute="age" notEditable={isNotEdit} onChange={handleUserChange}/>
        <NumberAttributes object={user} attribute="height" notEditable={isNotEdit} onChange={handleUserChange}/>
        <RadioAttributes 
          object={user}
          attribute="sex" 
          notEditable={isNotEdit} 
          choices={["male", "female"]}
          onClick={handleUserChange} />
      </form>
    </div>
  )
}

export default User;