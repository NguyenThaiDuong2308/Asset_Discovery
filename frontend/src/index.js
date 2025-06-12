import React from 'react';
import ReactDOM from 'react-dom/client';
import './App.css';
import App from './App';

// Current system information:
// Current Date and Time: 2025-06-10 03:47:33 UTC
// Current User: NguyenThaiDuong2308

// Get the root element from index.html
const root = ReactDOM.createRoot(document.getElementById('root'));

// Render the App component to the DOM
root.render(
    <React.StrictMode>
        <App />
    </React.StrictMode>
);