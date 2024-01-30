<script>
    import { onMount } from 'svelte';
    import dayjs from 'dayjs';
    import { writable, get } from 'svelte/store';
    import { startDate, endDate } from '$lib/stores';
    let start;
    let end;

    let availability = writable({});
    let tasks = writable([]);
    let taskInput = '';
    let assignedInput = '';

    onMount (() => {
        start = get(startDate);
        end = get(endDate);

        start = dayjs(start);
        end = dayjs(end);

        function initializeAvailability(start, end) {
            let days = {};

            for (let current = start; current.isBefore(end.add(1, 'day')); current = current.add(1, 'day')) {
                const dateString = current.format('YYYY-MM-DD');
                days[dateString] = new Array(24).fill(false);
            }
            availability.set(days);
        }
        initializeAvailability(dayjs(start), dayjs(end));
    });

    /** @param {number} day
     * @param {number} hour
     */
    function toggleSlot(day, hour) {
        availability.update(a => {
            a[day][hour] = !a[day][hour]
            return a;
        });
    }

    /** @param {string} taskDescription
     * @param {string} assigneeName
     */
    function addTask(taskDescription, assigneeName) {
        tasks.update(currentTasks => {
            const newTask = {
                id: currentTasks.length + 1,
                description: taskDescription,
                assignedTo: assigneeName,
                completed: false
            };
            return [...currentTasks, newTask]
        });
        taskInput = '';
        assignedInput = '';
    }

    /** @param {number} taskId */
    function toggleCompletion(taskId) {
        tasks.update(currentTasks => {
            const index = currentTasks.findIndex(t => t.id === taskId);
            if (index !== -1) {
                currentTasks[index].completed = !currentTasks[index].completed;
            }
            return currentTasks;
        });
    }

</script>

<header>

</header>

<main>
    <div class="content-wrap">
        <div class="menu-bar">
            <button class="menu-button">
                <img src="../menubar.png" alt="menu bar" class="hamburger-icon" >
                <span class="logo">LemmeKnow</span>
            </button>
            <button class="menu-button">
                <img src="../users.png" alt="menu bar" class="user-icon" >
                <span class="members-title">Members</span>
            </button>
        </div>
        <div>
            <img src="../chat.png" alt="chat" class="placeholder_images">
            <img src="../chat.png" alt="chat" class="placeholder_images">
        </div>
        <div class="calendar-container">
            <span class="calendar-title">AVAILABILITY CALENDAR</span>
            {#each Object.keys($availability) as day}
                <div class="day">
                    <h3>{day}</h3>
                    <div class="slots">
                        {#each $availability[day] as available, hour}
                            <div class="slot {available ? 'available' : ''}" on:click={() => toggleSlot(day, hour)}>
                                {hour}:00
                            </div>
                        {/each}
                    </div>
                </div>
            {/each}
            <form on:submit|preventDefault={() => addTask(taskInput, assignedInput)}>
                <input type="text" bind:value={taskInput} placeholder="Enter task description (50 characters max)" maxlength="50">
                <input type="text" bind:value={assignedInput} placeholder="Enter assignee name (50 characters max)" maxlength="50">
                <button type="submit" disabled={!taskInput.trim() || !assignedInput.trim()}>Add Task</button>
            </form>
            {#each $tasks as task (task.id)}
                <div class="task-item">
                    <input type="checkbox" bind:checked={task.completed} on:click={() => toggleCompletion(task.id)}>
                    <span class={task.completed ? 'completed-task' : ''}>{task.description}</span>
                    {#if task.assignedTo}
                        <span>Assigned to: {task.assignedTo}</span>
                    {/if}
                </div>
            {/each}
        </div>
    </div>
</main>

<style>
    placeholder_images {
        margin-top: 1rem;
    }

    .calendar-title {
        display: flex;
        align-items: flex-start;
        justify-content: center;
        flex-direction: column;
        text-align: center;
        font-size: 3rem;
        margin-top: 0.25rem;
        font-family: "Baloo Bhai 2";
        margin-left: 1rem;
        font-weight: bolder;
        color: black;
    }
    .menu-bar {
        position: relative;
        display: flex;
        flex-direction: column;
        top: 0;
        left: 0;
        align-items: flex-start;
    }

    .user-icon {
        width: 3rem;
        display: block;
        margin-left: 1.5rem;
    }

    .menu-button {
        background: none;
        border: none;
        cursor: pointer;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        line-height: 1;
        text-align: center;
    }

    .menu-button:focus {
        outline: none;
    }

    .logo {
        font-size: 1.5rem;
        margin-top: 0.25rem;
        font-family: "Baloo Bhai 2";
        font-weight: bolder;
        color: #73a0e7;
    }

    img.hamburger-icon {
        width: 3rem;
        margin: 0 auto;
        display: block;
    }

    .members-title {
        display: flex;
        align-items: flex-start;
        justify-content: center;
        text-align: center;
        font-size: 1.5rem;
        margin-top: 0.25rem;
        margin-left: 1rem;
        font-family: "Baloo Bhai 2";
        font-weight: bolder;
        color: black;
    }

    .menu-button:hover .hamburger-icon {
        transform: scale(1.2);
    }

    .menu-button:hover .user-icon {
        transform: scale(1.2);
    }

    .content-wrap {
        display: flex;
        flex-direction: row;
    }

    .calendar-container {
        display: flex;
        flex-direction: column;
        flex-wrap: wrap;
        margin-left: 4rem;
        margin-top: 3rem;
    }

    .day {
        border: 1px solid #ccc;
        padding: 10px;
    }

    .slots {
        display: grid;
        grid-template-columns: repeat(4, 1fr);
        gap: 5px;
    }

    .slot {
        padding: 5px;
        background-color: #f0f0f0;
        text-align: center;
        cursor: pointer;
    }

    .slot.available {
        background-color: #a9e1a9;
    }

    input {
        padding: 0.5rem;
        margin-bottom: 0.5rem;
        margin-top: 0.5rem;
        width: 80%;
        max-width: 15rem;
        text-align: center;
        font-size: 1rem;
        background-color: #c9e7e7;
        border-radius: 15px;
        border: 2px solid transparent;
    }

    button {
        padding: 0.5rem 1rem;
        background-color: #2774d0;
        color: white;
        border: none;
        align-items: center;
        font-size: 1rem;
        border-radius: 1rem;
        cursor: pointer;
    }

    button:hover {
        background-color: gray;
        color: white;
    }

    .task-item {
        display: flex;
        margin-bottom: 1rem;
        margin-top: 0.5rem;
        padding: 0.5rem;
        background-color: #f9f9f9;
        border-raius: 0.2rem;
        font-family: "Baloo Bhai 2";
        align-items: center;
        font-size: 1.25rem;
    }

    .task-item input[type="checkbox"] {
        accent-color: #879DB7;
        transform: scale(1.5);
        cursor: pointer;
    }

    .task-item .completed-task {
        text-decoration: line-through;
        color: #879DB7;
    }

    .task-item span {
        margin-left: 1rem;
        color: #333;
        font-weight: bold;
    }

    button[type="submit"]:disabled {
        background-color: #ccc;
        cursor: not-allowed;
    }
</style>