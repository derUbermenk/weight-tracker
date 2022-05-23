import React, {useEffect, useState} from "react";
import {RadioAttributes, TextAttributes, NumberAttributes, Button} from './UserComponents'

async function saveUser(user) {
  const requestUrl = 'http://localhost:8080/v1/api/user'
  const requestOptions = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(user)
  }

  const request = new Request(requestUrl, requestOptions)

  const response = await fetch(request);
  response = await response.json();

  // do something with the response here
  alert(response.status)
}

function NewUser() {
  const [user, setUser] = useState({});

  const handleUserChange = (attribute, value) => {
    var updated_user = user
    updated_user[attribute] = value
    setUser({...updated_user})
  }

  const handleSave = () => {
    saveUser(user)
  }

  return (
    <div>
      <h1>User Creation</h1>
      <div>
        {JSON.stringify(user)}
      </div>
      <form onSubmit={handleSave}>
        <TextAttributes object={user} attribute="name" notEditable={false} onChange={handleUserChange}/>
        <TextAttributes object={user} attribute="email" notEditable={false} onChange={handleUserChange}/>
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
          onClick={handleUserChange} />

        <br></br>
        <input type="submit" value="Create User"></input>

      </form>
    </div>
  )
}

export default NewUser; 