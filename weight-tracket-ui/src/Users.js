import React, {Component} from "react";

function Table(props) {
  const { data_values, data_headers } = props

  return(
    <table>
      <tr>
        {data_headers.map((data_header) => {
          return <th>{data_header}</th>
        })}
      </tr>
      
      {
        data_values.map( (entry) => {
          return(
            <tr>
              {data_headers.map((data_header) => {
                return <td>{entry[data_header]}</td>
              }
              )}
            </tr>
          )
        })
      }
    </table>
  )
}

class Users extends Component {
  constructor(props) {
    super(props);

    this.state = {
      users: [],
      keys: [],
    }

  }

  setup_data() {
    fetch("http://localhost:8080/v1/api/user", {mode: "cors"})
      .then(response => response.json())
      .then(data => this.setState({
        users: data,
        keys: Object.keys(data[0])
      }))
  }

  componentDidMount() {
    this.setup_data()
  }

  render() {
     const { users, keys } = this.state
     // const keys = Object.keys(this.state.users[0])

    return (
      <div>
        <div>
          <a href="/hello">Hello</a>
        </div>

        <h1>Users</h1>
        <Table data_values={users} data_headers={keys} />
      </div>
    )
  };
}

export default Users;
