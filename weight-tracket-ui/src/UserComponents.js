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

export {
  TextAttributes,
  NumberAttributes,
  RadioAttributes,
  Button
}