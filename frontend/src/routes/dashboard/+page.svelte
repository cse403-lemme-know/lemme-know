<script>
    let membersOnline = false;
    let chatMessages = [
        { user: "brooklyn", content: "running late" },
        { user: "zohar", content: "I hurt my foot"},
        { user: "finn", content: "just pushed more changes"}
    ];


    let poll = {
        question: "Where to go for lunch?",
        options: [
            { option: "KFC", votes: 1 },
            { option: "MCD's", votes: 2 },
            { option: "CHI MAC", votes: 7 }
        ],
        totalVotes: 10,
        active: true
    };

    let tasks = [];
    let availability = {};
    function toggleMembers() {
        membersOnline = !membersOnline;
    }
</script>

<header>
    <div class="menu-bar">
        <button on:click={toggleMembers} class="menu-button">
            <img src="../menubar.png" alt="menu bar" class="hamburger-icon" >
            <span class="logo">LemmeKnow</span>
        </button>
        <button on:click={toggleMembers} class="menu-button">
            <img src="../users.png" alt="menu bar" class="hamburger-icon" >
            <span class="members-title">Members</span>
        </button>

    </div>
</header>

<main>
    <div>

    </div>
</main>

<style>
    .menu-bar {
        position: relative;
        display: flex;
        flex-direction: column;
        top: 0;
        left: 0;
        align-items: flex-start;
    }

    .menu-button {
        background: none;
        border: none;
        cursor: pointer;
        display: inline-flex;
        flex-direction: column;
        padding: 0.625rem;
        line-height: 1;
    }

    .menu-button:focus {
        outline: none;
    }

    img {
        width: 20%;
        margin-left: 25%;
    }

    .hamburger-icon {
        width: 20%;
    }

    .logo {
        font-size: 1.5rem;
        margin-top: 0.25rem;
        font-family: "Baloo Bhai 2";
        font-weight: bolder;
        color: #73a0e7;
    }

    .members-title {
        text-align: center;
        font-size: 1.5rem;
        margin-top: 0.25rem;
        margin-left: 1rem;
        font-family: "Baloo Bhai 2";
        font-weight: bolder;
        color: black;
    }

    .menu-button:hover .hamburger-icon {
        background-color: #d3d3d3;
        border-radius: 0.5rem;
        transform: scale(1.2);
    }
</style>