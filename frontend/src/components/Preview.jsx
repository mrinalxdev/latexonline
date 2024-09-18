import React from 'react';
import { Document, Page } from 'react-pdf';

function Preview({ pdf }) {
  return (
    <div className="w-1/2 p-4">
      {pdf ? (
        <Document file={{ data: atob(pdf) }}>
          <Page pageNumber={1} />
        </Document>
      ) : (
        <div className="w-full h-full border rounded bg-white flex items-center justify-center">
          <p>Compiled PDF will appear here</p>
        </div>
      )}
    </div>
  );
}

export default Preview;