import React from 'react';
import { NavLink } from 'react-router-dom';

const Sidebar = () => {
    return (
        <aside className="app-sidebar">
            <nav>
                <ul>
                    <li>
                        <NavLink to="/dashboard" className={({ isActive }) => isActive ? 'active' : ''}>
                            Dashboard
                        </NavLink>
                    </li>
                    <li>
                        <NavLink to="/assets" className={({ isActive }) => isActive ? 'active' : ''}>
                            Assets
                        </NavLink>
                    </li>
                    <li>
                        <NavLink to="/logs" className={({ isActive }) => isActive ? 'active' : ''}>
                            Logs
                        </NavLink>
                    </li>
                </ul>
            </nav>
        </aside>
    );
};

export default Sidebar;