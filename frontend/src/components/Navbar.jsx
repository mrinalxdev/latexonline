import React from 'react';

function Navbar({ onCompile }) {
  return (
    <nav className="bg-gray-800 text-white p-4">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-bold">LaTeX IDE</h1>
        <button 
          onClick={onCompile}
          className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
        >
          Compile
        </button>
      </div>
    </nav>
  );
}

export default Navbar;