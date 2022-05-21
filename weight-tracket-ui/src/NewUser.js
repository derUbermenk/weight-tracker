import React, {useEffect, useState} from "react";
import {RadioAttributes, TextAttributes, NumberAttributes, Button} from './UserComponents'

function NewUser() {
  const [user, setUser] = useState({});

  const handleUserChange = (attribute, value) => {
    var updated_user = user
    updated_user[attribute] = value
    setUser({...updated_user})
  }

  const handleSave = () => {
    alert(`tried saving ${user}`)
  }

  /*
	Name          string `json:"name"`
	Age           int    `json:"age"`
	Height        int    `json:"height"`
	Sex           string `json:"sex"`
	ActivityLevel int    `json:"activity_level"`
	WeightGoal    string `json:"weight_goal"`
	Email         string `json:"email"
  */

  return (
    <div>
      <h1>User Creation</h1>
      <div>
        {JSON.stringify(user)}
      </div>
      <form >
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
        <TextAttributes object={user} attribute="activity_level" notEditable={false} onChange={handleUserChange}/>
        <TextAttributes object={user} attribute="weight_goal" notEditable={false} onChange={handleUserChange}/>

        <Button onclick={handleSave} name="Save User" />
      </form>
    </div>
  )
}

export default NewUser; 