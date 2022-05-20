import React, {Component, useEffect, useState} from "react";
import {useParams} from "react-router-dom";

function sentenceCase(string) {
  return string[0].toUpperCase() + string.slice(1)
}

function TextAttributes(props) {
  const { object, attribute, notEditable, onChange} = props

  const handleChange = (e) => {
    const attribute = e.target.getAttribute("id")
    const value = e.target.value
    onChange(attribute, value)
  }

  return(
    <div>
      <label htmlFor={attribute}>{sentenceCase(attribute)}: </label>
      <input type="text" id={attribute} value={object[attribute]} 
             disabled={notEditable} onChange={handleChange}
      />
    </div>
  )
}

function NumberAttributes(props) {
  const { object, attribute, notEditable, onChange } = props

  const handleChange = (e) => {
    const attribute = e.target.getAttribute("id")
    const value = e.target.value
    onChange(attribute, value)
  }

  return(
    <div>
      <label htmlFor={attribute}>{sentenceCase(attribute)}</label>
      <input type="number" id={attribute} value={object[attribute]}
             disabled={notEditable} onChange={handleChange}
      />
    </div>
  )
}

function RadioAttributes(props) {
  const { object, attribute, notEditable, choices, onClick} = props

  const handleClick = (e) => {
    const value = e.target.value
    const attribute = e.target.getAttribute("name")

    onClick(attribute, value)
  }

  return(
    <div>
      <label htmlFor={attribute}>{sentenceCase(attribute)}: </label>
      {
        choices.map((choice) => {
          return (
            <span key={choice}>
              <input type="radio" id={`${attribute}_${choice}`}  value={choice}
                      name={attribute} disabled={notEditable}
                      checked={choice==object[attribute] ? true : false }
                      onClick={handleClick}
                      />
              <label htmlFor={choice}>{sentenceCase(choice)}</label>
            </span>
          )
        })
      }
    </div>
  )
}

function Button(props) {
  const { onclick, name } = props

  return(
    <button onClick={onclick}>
      {name}
    </button>
  )
}

function User() {
  var {userId} = useParams(); 
  const [user, setUser] = useState({});

  useEffect(() => { 
    fetch(`http://localhost:8080/v1/api/user/${userId}`, {mode: "cors"})
    .then(response => response.json())
      .then(user => setUser(user))
    
  }, []);

  const handleUserChange = (attribute, value) => {
    var updated_user = user
    updated_user[attribute] = value
    setUser({...updated_user})
  }

  const [isNotEdit, setisNotEdit] = useState(true);

  const handleEdit = () => {
    setisNotEdit(false) // make user editable
  }

  const handleSave = () => {
    alert(`Save Button clicked, saving ${JSON.stringify(user)}`)
    // saveUser
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