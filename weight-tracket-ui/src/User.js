import React, {Component, useEffect, useState} from "react";
import {useParams} from "react-router-dom";

function User() {
  var {userId} = useParams(); 
  const [user, setUser] = useState({});
  useEffect(() => { 
    fetch(`http://localhost:8080/v1/api/user/${userId}`, {mode: "cors"})
    .then(response => response.json())
      .then(user => setUser(user))
  })

  return (
    <div>
      <h1>Hello {`${user.name}`}!</h1>
      <br></br>
      <form>
        <div>
          <label for="name">Name: </label>
          <input type="text" id="name" value={user.name} disabled={true}/>
        </div>

        <div>
          <label for="email">Email: </label>
          <input type="text" id="email" value={user.email} disabled={true}/>
        </div>

        <div>
          <label for="age">Age: </label>
          <input type="number" id="age" value={user.age} disabled={true}/>
        </div>

        <div>
          <label for="height">Height: </label>
          <input type="number" id="height" value={user.height} disabled={true}/>
        </div>

        <div>
          <label for="sex">Sex: </label><br></br>
          <input type="radio" id="male"  value="male" 
                 name="sex" checked={user.sex=="male"} 
                 disabled={true} />
          <label for="male">Male</label>
          <br></br>
          <input type="radio" id="female"  value="female" 
                 name="sex" checked={user.sex=="female"} disabled={true} />
          <label for="female">Female</label>
        </div>
      </form>
    </div>
  )
}

export default User;