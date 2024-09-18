import React from 'react';

function StatusBar({ status }) {
  return (
    <div className="bg-gray-300 p-2 text-sm">
      {status}
    </div>
  );
}

export default StatusBar;