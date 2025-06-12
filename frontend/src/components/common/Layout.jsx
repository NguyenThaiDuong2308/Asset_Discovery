import React from 'react';
import Header from './Header';
import Sidebar from './Sidebar';

const Layout = ({ children }) => {
    return (
        <div className="app-container">
            <Header />
            <div className="main-content">
                <Sidebar />
                <div className="content-area">
                    {children}
                </div>
            </div>
        </div>
    );
};

export default Layout;