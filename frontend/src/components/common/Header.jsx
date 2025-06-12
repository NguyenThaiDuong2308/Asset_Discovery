import React from 'react';
import { Link } from 'react-router-dom';

const Header = () => {
    return (
        <header className="app-header">
            <div className="logo">
                <Link to="/">Asset Management System</Link>
            </div>
            <div className="user-info">
                <span>NguyenThaiDuong2308</span>
            </div>
        </header>
    );
};

export default Header;