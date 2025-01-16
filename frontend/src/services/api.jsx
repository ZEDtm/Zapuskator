const getServerInfo = async (url) => {
    try {
        // Отправляем запрос на сервер
        const endpoint = "/resto/get_server_info.jsp?encoding=UTF-8";
        const fullUrl = `${url}${endpoint}`;

        const response = await fetch(fullUrl);
        if (!response.ok) {
            console.log(`HTTP error! status: ${response.status}`);
        }

        const xmlData = await response.text();

        // Парсим XML
        const parser = new DOMParser();
        const xmlDoc = parser.parseFromString(xmlData, "text/xml");

        const version = xmlDoc.getElementsByTagName("version")[0].textContent;
        const state = xmlDoc.getElementsByTagName("serverState")[0].textContent;
        const edition = xmlDoc.getElementsByTagName("edition")[0].textContent;

        return ({ edition, version, state });
    } catch (err) {
        console.error(err);
    }
};

export default getServerInfo;