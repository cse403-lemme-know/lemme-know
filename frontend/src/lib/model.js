async function getUser() {
    try {
        const response = await fetch(`//${location.host}/api/user/`);
        const user = await response.json();
        return user;
    } catch (e) {
        return null;
    }
}

async function createGroup(name) {
    try {
        const response = await fetch(`//${location.host}/api/group/`, {
            method: "PATCH",
            body: JSON.stringify({name})
        });
        const result = await response.json();
        return result.groupId;
    } catch (e) {
        return null;
    }
}

getUser().then(console.log);

export {createGroup};