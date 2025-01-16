import React from 'react';

const ServersList = ({ servers, setServers, setUrl, setInput}) => {
    const changeUrl = (url) => {
        setUrl(url);
    };

    const deleteUrl = (url) => {
        const currentServers = servers.filter((server) => server !== url);
        setServers(currentServers);
    };

    return (
        <div className="servers-list">
            <h2>История</h2>
            {servers.map((server) => (
                <div key={server} className="server-item">
                    <button onClick={() => changeUrl(server)} className="server-button">
                        {server}
                    </button>
                    <button onClick={() => deleteUrl(server)} className="delete-button">
                        X
                    </button>
                </div>
            ))}
        </div>
    );
};

export default ServersList;