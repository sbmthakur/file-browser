import { useState, useEffect } from 'react';
import { useLocation, useParams } from 'react-router-dom';
import './App.css';

interface FilesDetails {
  name: string;
  size: number;
  type: string;
}

interface Resp extends FilesDetails {
  contents: FilesDetails[]
}

function App() {
  //const [fileData, setFileData] = useState({ name: 'some' })
  const fileStateData = {
    name: '',
    size: 0,
    type: 'file',
    contents: []
  }

  const { dirPath } = useParams()

  //let { pathname } = useLocation()
  //alert(pathname)
  //pathname = pathname.replace("/folders", '')
  //const initialPath = pathname === '/' ? '/' : pathname.substring(1, pathname.length)
  alert(dirPath)
  const initialPath = dirPath ? dirPath : '/'

  const [fileData, setFileData] = useState<Resp>(fileStateData)
  const [path, setPath] = useState(initialPath)

  useEffect(() => {
    async function getData() {
      const response = await fetch("http://localhost:1323/files?path=" + path)
      const data = await response.json()
      setFileData(data)
    }
    getData()
  }, [path])

  return (
    <div className="App">
      <header className="App-header">
        <ul id="Bread-crumb" style={{ listStyleType: 'none' }}>
          <li>Home</li>
          {
            path.split('/').map(item => {
              if(item) {
                return <li> {'>'} {item}</li>
              } else {
                return null
              }
              
            })
          }
        </ul>
        <table>
          <th>Name</th>
          <th>Size</th>
          { 
          fileData.contents.map(o => {
            return <tr><td><span className="File-name" onClick={() => { 
              if (o.type === "dir") {
                const newName = path === '/' ? o.name : path + '/' + o.name
                setPath(newName)
              }
            }}>{o.name}</span></td><td>{o.type === "dir" ? 0: o.size}</td></tr>
          })
          }
        </table>
      </header>
    </div>
  );
}

export default App;
