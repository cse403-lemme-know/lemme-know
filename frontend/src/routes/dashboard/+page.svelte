<script>
    let availability = {};
    let membersOnline = false;
    let chatMessages = [
        { user: "brooklyn", content: "running late" },
        { user: "zohar", content: "I hurt my foot"},
        { user: "finn", content: "just pushed more changes"}
    ];
    let tasks = [];
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
    </div>
</header>

<main>
    <div>

    </div>
</main>

<style>
    .menu-bar {
        position: absolute;
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

    hamburger-icon {
        width: 100%;
    }

    .logo {
        font-size: 1.5rem;
        margin-top: 0.25rem;
        font-family: "Baloo Bhai 2";
        font-weight: bolder;
        color: #73a0e7;
    }

    .menu-button:hover .hamburger-icon {
        background-color: #d3d3d3;
        border-radius: 0.5rem;
        transform: scale(1.2);
    }
</style>