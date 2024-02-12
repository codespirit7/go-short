import { useState } from 'react'
import { CopyToClipboard } from 'react-copy-to-clipboard';
import './App.css'

function App() {

  const [input, setInput] = useState("")
  const [url, setUrl] = useState("")
  const [copyState, setCopyState] = useState("Copy url")
  const [copyStatus, setCopyStatus] = useState(false);

  const handleUrl = async () => {
    try{
      const response = await fetch("http://localhost:8080/short-url", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          url: input
        })
      })

      const data = await response.json();
      console.log(data["short-url"]);
      setUrl(data["short-url"])
    }catch(err){
      console.err(err);
    }
  }

  const onCopyText = () => {
    setCopyStatus(true);
    setCopyState("Copied url")
    setInput("")
    setTimeout(() => setCopyState("Copy url"), 2000); // Reset status after 2 seconds
  };
  return (
    <>
        <section className="wrapper">
          <div className="top">GoShort</div>
          <div className="bottom" aria-hidden="true">GoShort</div>
        </section>
        <span className='humans-url'> Short url for Humans & Aliens</span>
        <div>
        <input value = {input} onChange={(e) => setInput(e.target.value)} placeholder='Enter URL'/>
          <button  onClick={handleUrl}>Go Short</button>
        </div>
        <div className='main'>
          <div className='url'>
            {url}
          </div>
          <div>
          {url && <CopyToClipboard text={url} onCopy={onCopyText}>
           <button>{copyState}</button>
           </CopyToClipboard>}
          </div>
          <div>
        {url && <iframe
          className='iframe'
          width={900}
          height={900}
          src={url}>
    </iframe>}
        </div>
        </div>
        
    </>
  )
}

export default App
