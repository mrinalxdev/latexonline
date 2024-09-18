import React from 'react';
import Navbar from './components/Navbar';
import Editor from './components/Editor';
import Preview from './components/Preview';
import StatusBar from './components/StatusBar';

function App() {
  const [latex, setLatex] = React.useState('');
  const [compiledPdf, setCompiledPdf] = React.useState(null);
  const [status, setStatus] = React.useState('');

  const compileLaTeX = async () => {
    setStatus('Compiling...');
    try {
      const response = await fetch('/compile', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ latex }),
      });
      const result = await response.json();
      if (result.error) throw new Error(result.error);
      setCompiledPdf(result.pdf);
      setStatus('Compilation successful!');
    } catch (error) {
      setStatus(`Error: ${error.message}`);
    }
  };

  return (
    <div className="h-screen flex flex-col bg-gray-100">
      <Navbar onCompile={compileLaTeX} />
      <div className="flex-grow flex">
        <Editor latex={latex} setLatex={setLatex} />
        <Preview pdf={compiledPdf} />
      </div>
      <StatusBar status={status} />
    </div>
  );
}

export default App;