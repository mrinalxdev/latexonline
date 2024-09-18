import React, { useState } from "react";
import AceEditor from "react-ace";
import "ace-builds/src-noconflict/mode-latex";
import "ace-builds/src-noconflict/theme-monokai";
import "ace-builds/src-noconflict/ext-language_tools";

function Editor({ latex, setLatex }) {
  const [status, setStatus] = useState("");
  const compileLaTeX = async () => {
    setStatus("Compiling...");
    try {
      const response = await fetch("/compile", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ latex }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const result = await response.json();

      if (result.error) {
        throw new Error(result.error);
      }

      // Handle successful compilation
      setStatus("Compilation successful!");
    } catch (error) {
      console.error("Compilation error:", error);
      setStatus(`Error: ${error.message}`);
    }
  };
  return (
    <div className="w-1/2 p-4">
      <AceEditor
        mode="latex"
        theme="monokai"
        onChange={setLatex}
        value={latex}
        name="UNIQUE_ID_OF_DIV"
        editorProps={{ $blockScrolling: true }}
        setOptions={{
          enableBasicAutocompletion: true,
          enableLiveAutocompletion: true,
          enableSnippets: true,
        }}
        style={{ width: "100%", height: "100%" }}
      />
      <button onClick={compileLaTeX}>Compile</button>
      <p>{status}</p>
    </div>
  );
}

export default Editor;
