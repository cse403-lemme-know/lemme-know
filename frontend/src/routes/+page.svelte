<script>
    import { goto } from '$app/navigation';
    function handleButtonClick() {
        goto("/dashboard");
    }
    let name = 'LemmeKnow';
</script>

<header>
    <h1>{name}</h1>
    <nav>
        <a href="/features/">FEATURES</a>
        <a href="/privacy">PRIVACY</a>
        <a href="/devteam">DEV TEAM</a>
    </nav>
</header>

<main>
    <h3>SCHEDULE HANGOUTS, PLAN ROAD-TRIPS, SHARE CALENDARS, EVERYTHING, EVERYWHERE, ALL AT ONCE.</h3>
    <div class="images">
        <div class="container">
            <img src="/highfive.png" alt="3 people trying to figure out their collective availability">
            <img src="/cal.png" alt="family of four sharing their calendars">

        </div>
        <div class="input-container">
            <input type="text" placeholder="Enter Group name..">
            <button on:click={handleButtonClick}>Let me know!</button>
        </div>
        <div class="container">
            <img src="/road.png" alt="group of friends driving in the car">
            <img src="/calendar.png" alt="three colleagues figuring out their work schedule">
        </div>
    </div>
</main>

<style>
    main {
        max-width: 1200px;
        margin: 0 auto;
        padding: 2rem;
        display: flex;
        flex-direction: column;
        align-items: center;
        text-align: center;
        font-family: "Baloo Da 2", sans-serif;
    }

    header {
        width: 100%;
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1rem 0;
    }

    nav a {
        text-decoration: none;
        color: #000;
        margin: 0 1rem;
        font-weight: bold;
        font-family: "Baloo Da 2"
    }

    h1 {
        margin: 0.5rem 0;
        font-family: "Baloo Bhai 2";
        padding: 0 2rem;
    }

    h3 {
        margin: 0.5rem 0;
        font-family: "Baloo Da 2";
        font-weight: bolder;
        font-size: 2rem;
        color: #5D75CD;
    }

    img {
        width: 100%;
        max-width: 500px;
        margin: 0.5rem 0.5rem;
    }

    .images {
        display: flex;
        justify-content: space-around;
        align-items: center;
        flex-wrap: wrap;
        width: 100%;
        gap: 1rem;
    }

    .container {
        flex: 1;
        max-width: 50%;
        margin: 1rem;
    }

    .input-container {
        display: flex;
        flex: 1;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        margin: 1rem;
        border-color: white;
        color: white;
    }

    input {
        padding: 0.5rem;
        margin-bottom: 0.5rem;
        width: 80%;
        max-width: 300px;
        text-align: center;
        font-size: 1rem;
        background-color: #EEDAF4;
        border-radius: 15px;
        border: 2px solid transparent;
    }

    input:focus {
        outline: none;
        border-color: black;
    }

    button {
        padding: 0.5rem 1rem;
        background-color: #879DB7;
        color: white;
        border: none;
        font-size: 1rem;
        border-radius: 1rem;
        cursor: pointer;
    }

    button:hover {
        background-color: gray;
        color: white;
    }

    @media (min-width: 640px) {
        .images, .container, .input-container {
            max-width: 100%;
        }

        nav a {
            margin: 0.5rem;
        }
    }
</style>