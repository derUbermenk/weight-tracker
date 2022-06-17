import React, {useEffect, useState} from "react";
import  { useNavigate } from 'react-router-dom'
import {RadioAttributes, TextAttributes, NumberAttributes, Button} from './UserComponents'

/**
 * 
 * @param {*} user 
 * @returns {Promise} response 
 */
async function saveUser(user) {
  const requestUrl = 'http://localhost:8080/v1/api/user'
  const requestOptions = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(user)
  }

  const request = new Request(requestUrl, requestOptions)

  const response = await fetch(request);
  const json = await response.json()
  return json 
}

function NewUser() {
  const [user, setUser] = useState({});
  var navigate = useNavigate();

  const handleUserChange = (attribute, value) => {
    var updated_user = user
    updated_user[attribute] = value
    setUser({...updated_user})
  }

  const handleSave = async (e) => {
    e.preventDefault(); 

    // send save user request; then get status
    const json = await saveUser(user)
    const status = json['Status']
    const data = json['Data']
    const user_id = json['UserID']

    // do something about status
    if (status == 'success') {
      navigate(`/user/${user_id}`)
    } else {
      alert(`${status} because ${data}`)
      // do not change user
    }
  }

  return (
    <div>
      {
      }
      <h1>User Creation</h1>
      <div>
        {JSON.stringify(user)}
      </div>
      <form onSubmit={handleSave}>
        <TextAttributes object={user} attribute="name" notEditable={false} onChange={handleUserChange} type="text"/>
        <TextAttributes object={user} attribute="email" notEditable={false} onChange={handleUserChange} type="email"/>
        <NumberAttributes object={user} attribute="age" notEditable={false} onChange={handleUserChange}/>

        <RadioAttributes 
          object={user}
          attribute="sex" 
          notEditable={false} 
          choices={["male", "female"]}
          onClick={handleUserChange} />

        <br></br>

        <NumberAttributes object={user} attribute="height" notEditable={false} onChange={handleUserChange}/>
        <NumberAttributes object={user} attribute="activity_level" notEditable={false} onChange={handleUserChange}/>
        <RadioAttributes 
          object={user}
          attribute="weight_goal" 
          notEditable={false} 
          choices={["loose", "maintain", "gain"]}
          onClick={handleUserChange}
        />

        <br></br>

        <input type="submit" value="Create User"></input>

      </form>
    </div>
  )
}

export default NewUser; 