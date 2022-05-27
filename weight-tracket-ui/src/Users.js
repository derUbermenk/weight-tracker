import React, {useEffect, useState} from "react";

function Table(props) {
  const { data_values, data_headers, onDelete } = props

  const handleUserDelete = (e) => {
    const userID = e.target.getAttribute("id")
    onDelete(userID)
  } 

  return(
    <table>
      <thead>
        <tr>
          {data_headers.map((data_header) => {
            return <th key={data_header}>{data_header}</th>
          })}
        </tr>
      </thead>
      
      <tbody>
        {
          data_values.map( (entry) => {
            return(
              <tr key={entry.id}>
                {data_headers.map((data_header) => {
                  return <td key={data_header}>{entry[data_header]}</td>
                }
                )}
                <td><a href={`/user/${entry.id}`}>{entry.name}</a></td>
                <td>
                  <span id={entry.id} onClick={handleUserDelete} className="material-icons md-18">close</span>
                </td>
              </tr>
            )
          })
        }
      </tbody>
    </table>
  )
}

async function getUsers(usersSetter, keysSetter) {
  const requestUrl = `http://localhost:8080/v1/api/user`
  const requestOptions = {
    mode: 'cors'
  }

  const request = new Request(requestUrl, requestOptions)

  const response = await fetch(request);
  const users = await response.json();

  usersSetter(users)

  // sample a user
  const sample_user = users[0]
  const keys = Object.keys(sample_user)
  keysSetter(keys)
}

function getKeys(users, keySetter) {
  // get first user value
  const user_sample = Object.values(users)
  const keys = Object.keys(user_sample)

  keySetter(keys)
}

async function deleteUser(userID) {
  const requestUrl = `http://localhost:8080/v1/api/user/${userID}`
  const requestOptions = {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
  }
  const request = new Request(requestUrl, requestOptions) 

  const response = await fetch(request);
  const json = await response.json()
  return json
}


function Users() {
  const [users, setUsers] = useState({});
  const [keys, setKeys] = useState([]);

  // fetch users at render
  useEffect(() => {
    getUsers((users) => setUsers(users), (keys) => setKeys(keys));
    },
    []
  )

  const handleUserDeletion = async (userID) => {
    const json = await deleteUser(userID)
    const status = json['Status']
    const data = json['Data']

    if (status == 'success') {
      const deletedUserID = json['UserID']
      const updated_user_dict = delete users[deletedUserID]
      setUsers(updated_user_dict)
    } else {
      alert(`${status} because ${data}`)
    }
  } 

  return (
    <div>
      <div>
      </div>

      <h1>Users</h1>
      <Table data_values={Object.values(users)} data_headers={keys} onDelete={handleUserDeletion} />

      <br></br>
      <a href="/user/new">Add User</a>
    </div>
  )
}

export default Users;
