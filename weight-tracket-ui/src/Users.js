import React, {Component} from "react";

function Table(props) {
  const { data_values, data_headers } = props

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
                  <span onClick={() => alert("tried delete")} className="material-icons md-18">close</span>
                </td>
              </tr>
            )
          })
        }
      </tbody>
    </table>
  )
}

class Users extends Component {
  constructor(props) {
    super(props);

    this.state = {
      users: {},
      keys: []
    }

  }

  setup_data() {
    fetch("http://localhost:8080/v1/api/user", {mode: "cors"})
      .then(response => response.json())
      .then(data => {
        var fetched_users = {}
        for(const user of data) {
          fetched_users[user.id] = user
        }

        this.setState({
          users: fetched_users,
          keys: Object.keys(data[0])
        })
      })
  }

  componentDidMount() {
    this.setup_data()
  }

  render() {
     const { users, keys } = this.state

    return (
      <div>
        <div>
        </div>

        <h1>Users</h1>
        <Table data_values={Object.values(users)} data_headers={keys} />

        <br></br>
        <a href="/user/new">Add User</a>
      </div>
    )
  };
}

export default Users;
