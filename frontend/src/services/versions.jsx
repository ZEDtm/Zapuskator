const getPathOnVersion = (edition, version) => {
    const versionParts = version.split('.');
    const versionCurrent = versionParts[0][0] + versionParts[1][0] + versionParts[2][0]
    const backOfficesPath = {
        "741": "Q:\\BackOffice\\RMS\\Office741\\BackOffice.exe",
        "742": "Q:\\BackOffice\\RMS\\Office742\\BackOffice.exe",
        "746": "Q:\\BackOffice\\RMS\\Office746\\BackOffice.exe",
        "747": "Q:\\BackOffice\\RMS\\Office747\\BackOffice.exe",
        "756": "Q:\\BackOffice\\RMS\\Office756\\BackOffice.exe",
        "766": "Q:\\BackOffice\\RMS\\Office766\\BackOffice.exe",
        "767": "Q:\\BackOffice\\RMS\\Office767\\BackOffice.exe",
        "772": "Q:\\BackOffice\\RMS\\Office772\\BackOffice.exe",
        "774": "Q:\\BackOffice\\RMS\\Office774\\BackOffice.exe",
        "776": "Q:\\BackOffice\\RMS\\Office776\\BackOffice.exe",
        "777": "Q:\\BackOffice\\RMS\\Office777.009\Office\\BackOffice.exe",
        "778": "Q:\\BackOffice\\RMS\\Office778\\BackOffice.exe",
        "786": "Q:\\BackOffice\\RMS\\Office786\\BackOffice.exe",
        "791": "Q:\\BackOffice\\RMS\\Office791\\BackOffice.exe",
        "796": "Q:\\BackOffice\\RMS\\Office796\\BackOffice.exe",
        "797": "Q:\\BackOffice\\RMS\\Office797\\BackOffice.exe",
        "806": "Q:\\BackOffice\\RMS\\\Office806\\BackOffice.exe",
        "816": "Q:\\BackOffice\\RMS\\Office816\\BackOffice.exe",
        "817": "Q:\\BackOffice\\RMS\\Office817\\BackOffice.exe",
        "826": "Q:\\BackOffice\\RMS\\Office826\\Office826\\BackOffice.exe",
        "827": "C:\\iiko_Distr\Office\RMS\Office827\\BackOffice.exe",
        "836": "Q:\\BackOffice\\RMS\\Office836\\BackOffice.exe",
        "837": "Q:\\BackOffice\\RMS\\Office837\\BackOffice.exe",
        "846": "Q:\\BackOffice\\RMS\\Office846\\BackOffice.exe",
        "847": "Q:\\BackOffice\\RMS\\Office847\\BackOffice.exe",
        "857": "Q:\\BackOffice\\RMS\\Office857\\BackOffice.exe",
        "858": "Q:\\BackOffice\\RMS\\Office858\\BackOffice.exe",
        "866": "Q:\\BackOffice\\RMS\\Office866\\BackOffice.exe",
        "867": "Q:\\BackOffice\\RMS\\Office867\\BackOffice.exe",
        "868": "Q:\\BackOffice\\RMS\\Office868\\Office\\BackOffice.exe",
        "869": "Q:\\BackOffice\\RMS\\Office869\\BackOffice.exe",
        "876": "Q:\\BackOffice\\RMS\\Office876\\BackOffice.exe",
        "877": "Q:\\BackOffice\\RMS\\Office877\\BackOffice.exe",
        "886": "Q:\\BackOffice\\RMS\\Office886\\BackOffice.exe",
        "887": "Q:\\BackOffice\\RMS\\Office887\\Office887\\BackOffice.exe",
        "888": "Q:\\BackOffice\\RMS\\Office888\\BackOffice.exe",
        "889": "Q:\\BackOffice\\RMS\\Office889\\BackOffice.exe",
        "896": "Q:\\BackOffice\\RMS\\Office896\\BackOffice.exe",
        "897": "Q:\\BackOffice\\RMS\\Office897\\Office897\\iikoOffice\\BackOffice.exe",
        "898": "C:\\BackOffice\\RMS\\Office898\\BackOffice.exe",
        "899": "Q:\\BackOffice\\RMS\\Office899\\BackOffice.exe",
        "906": "Q:\\BackOffice\\RMS\\Office906\\BackOffice.exe",
        "907": "C:\\BackOffice\\RMS\\Office907\\BackOffice.exe",
        "908": "C:\\BackOffice\\RMS\\Office908\\BackOffice.exe",
    };
    const chainsPath = {
        "736": "Q:\\BackOffice\\Chain\\COffice736\\BackOffice.exe",
        "756": "Q:\\BackOffice\\Chain\\COffice756\\BackOffice.exe",
        "778": "Q:\\BackOffice\\Chain\\COffice778\\BackOffice.exe",
        "786": "Q:\\BackOffice\\Chain\\Coffice786\\BackOffice.exe",
        "791": "Q:\\BackOffice\\Chain\\Coffice791\\BackOffice.exe",
        "797": "Q:\\BackOffice\\Chain\\COffice797\\BackOffice.exe",
        "806": "Q:\\BackOffice\\Chain\\Coffice806\\BackOffice.exe",
        "816": "Q:\\BackOffice\\Chain\\COffice816\\BackOffice.exe",
        "817": "Q:\\BackOffice\\Chain\\COffice817\\BackOffice.exe",
        "826": "Q:\\BackOffice\\Chain\\СOffice826\\BackOffice.exe",
        "827": "Q:\\BackOffice\\Chain\\СOffice827\\BackOffice.exe",
        "836": "Q:\\BackOffice\\Chain\\СOffice836\\BackOffice.exe",
        "837": "Q:\\BackOffice\\Chain\\СOffice837\\BackOffice.exe",
        "846": "Q:\\BackOffice\\Chain\\СOffice846\\BackOffice.exe",
        "857": "Q:\\BackOffice\\Chain\\СOffice857\\BackOffice.exe",
        "858": "Q:\\BackOffice\\Chain\\COffice858\\BackOffice.exe",
        "867": "Q:\\BackOffice\\Chain\\СOffice867\\BackOffice.exe",
        "868": "Q:\\BackOffice\\Chain\\СOffice868\\BackOffice.exe",
        "869": "C:\\iiko_Distr\Office\Chain\COffice869\\BackOffice.exe",
        "876": "Q:\\BackOffice\\Chain\\СOffice876\\Office\\BackOffice.exe",
        "877": "Q:\\BackOffice\\Chain\\СOffice877\\BackOffice.exe",
        "886": "Q:\\BackOffice\\Chain\\СOffice886\\BackOffice.exe",
        "888": "Q:\\BackOffice\\Chain\\СOffice888\\BackOffice.exe",
        "889": "Q:\\BackOffice\\Chain\\СOffice889\\BackOffice.exe",
        "897": "Q:\\BackOffice\\Chain\\СOffice897\\Office\\Office\\BackOffice.exe",
        "898": "Q:\\BackOffice\\Chain\\СOffice898\\BackOffice.exe",
        "907": "Q:\\BackOffice\\Chain\\СOffice906\\BackOffice.exe",
    };

    if (edition === 'default') {
        return backOfficesPath[versionCurrent];
    } else if (edition === 'chain') {
        return chainsPath[versionCurrent];
    } else {
        return '';
    }
};

export default getPathOnVersion;