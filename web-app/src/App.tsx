import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
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
  const fileStateData = {
    name: '',
    size: 0,
    type: 'file',
    contents: []
  }

  const { dirPath } = useParams()
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
          <li><a href={document.location.origin +'/'}>Home</a></li>
          {(() => {
            if (path === '/') {
              return 
            }
            const paths = path.split('/')
            const c = []
            for(let i = 0; i < paths.length; i++) {
              c.push(<li>{'> '}<a href={document.location.origin +'/' + paths.slice(0, i + 1).join('/')}>{paths[i]}</a></li>) 
            }
            return c
          })()}
        </ul>
        <table>
          <th>Name</th>
          <th>Size</th>
          { 
            fileData.contents.map(o => {
              return <tr><td>
                {
                  (() => {
                    const fileExtension = o.name.substring(o.name.lastIndexOf('.') + 1) 
                    const className = "file-icon fiv-viv fiv-icon-" + fileExtension
                    return <span className={className}></span>
                  })()
                }
                <span className="File-name" onClick={() => {
                  if (o.type === "dir") {
                    const newName = path === '/' ? o.name : path + '/' + o.name
                    setPath(newName)
                  }
                }}>{o.name}</span></td><td>{o.type === "dir" ? '-' : o.size}</td></tr>
            })
          }
        </table>
      </header>
    </div>
  );
}

export default App;
