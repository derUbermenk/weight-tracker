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
    // saveUser(user)
  }

  return (
    <div>
      <h1>User Creation</h1>
      <div>
        {JSON.stringify(user)}
      </div>
      <form>
        <TextAttributes object={user} attribute="name" notEditable={false} onChange={handleUserChange}/>
        <TextAttributes object={user} attribute="email" notEditable={false} onChange={handleUserChange}/>
        <NumberAttributes object={user} attribute="age" notEditable={false} onChange={handleUserChange}/>
        <NumberAttributes object={user} attribute="height" notEditable={false} onChange={handleUserChange}/>
        <RadioAttributes 
          object={user}
          attribute="sex" 
          notEditable={false} 
          choices={["male", "female"]}
          onClick={handleUserChange} />
      </form>
    </div>
  )
}

export default NewUser; 