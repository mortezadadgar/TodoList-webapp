"use strict";

// create todos element earlier
function createTodosElement() {
    const todos = document.createElement("div");
    todos.id = "todos";
    document.body.append(todos);
    return todos;
}

function createBasicUI(todos) {
    const input = document.createElement("input");
    input.id = "todo-input";
    input.type = "text";
    input.name = "todo-input";
    input.onkeydown = function (ev) {
        if (ev.key === "Enter") {
            createTodo(todos);
        }
    };

    const button = document.createElement("a");
    button.type = "submit";
    button.className = "todo-add";
    button.onclick = function () {
        createTodo(todos);
    };
    button.innerHTML = "<i class='fa-solid fa-square-plus fa-2xl'></i>";
    todos.append(input, button);
}

//TODO: Implement UI to edit todos
function createTodoUI(id, todos, text, isDone) {
    const todoText = document.createElement("label");
    todoText.className = "todo-text";
    todoText.innerText = text;

    const todoCheck = document.createElement("input");
    todoCheck.className = "todo-check";
    todoCheck.type = "checkbox";
    todoCheck.checked = isDone;
    todoCheck.name = id;
    todoCheck.onclick = function () {
        editTodo(this, todoText);
    };

    const todoActions = document.createElement("div");
    todoActions.className = "todo-actions";
    todoActions.hidden = true;

    const todoEdit = document.createElement("a");
    todoEdit.type = "button";
    todoEdit.className = "todo-edit";
    todoEdit.innerHTML = "<i class='fa-solid fa-pen'></i>";

    const todoRemove = document.createElement("a");
    todoRemove.type = "button";
    todoRemove.className = "todo-remove";
    todoRemove.innerHTML = "<i class='fa-solid fa-trash'></i>";
    todoRemove.onclick = function () {
        deleteTodo(todo, todoCheck);
    };

    const todo = document.createElement("div");
    todo.id = "todo";
    todo.onmouseover = function () {
        if (todoActions) {
            todoActions.hidden = false;
        }
    };
    todo.onmouseleave = function () {
        if (todoActions) {
            todoActions.hidden = true;
        }
    };

    todoActions.append(todoEdit, todoRemove);

    todo.append(todoCheck, todoText, todoActions);

    todos.append(todo);
}

async function getTodos(todosElement) {
    const resp = await fetch("/todo/list", {
        method: "GET",
    });
    if (!resp.ok) {
        throw "The response was not OK";
    }
    const data = await resp.json();
    const todos = data.data;
    for (let i = 0; i < todos.length; i++) {
        createTodoUI(
            todos[i].id,
            todosElement,
            todos[i].text,
            todos[i].is_done
        );
    }
    return data;
}

async function getLastID() {
    const resp = await fetch("/todo/id", {
        method: "GET",
    });
    if (!resp.ok) {
        throw "The response was not OK";
    }
    const data = await resp.text();
    return data;
}

async function sendTodos(input) {
    await fetch("/todo", {
        body: input.value,
        method: "POST",
    });
}

async function deleteTodo(todo, todoCheck) {
    todo.remove();

    console.log(`DELETE id: ${todoCheck.name}`);
    const encodedTodo = JSON.stringify({
        data: { id: parseInt(todoCheck.name) },
    });

    await fetch("/todo", {
        body: encodedTodo,
        method: "DELETE",
    });
}

//DONE: send todo data as json
//NOTE: we dont want to send multiple items right?
async function editTodo(todoCheck, todoText) {
    const text = todoText.innerText;
    const isDone = todoCheck.checked;
    console.log(`id: ${todoCheck.name}, text: ${text}, isDone: ${isDone}`);
    const encodedTodo = JSON.stringify({
        data: { id: parseInt(todoCheck.name), text: text, is_done: isDone },
    });

    await fetch("/todo", {
        body: encodedTodo,
        method: "PUT",
    });
}

function createTodo(todosElement) {
    const input = document.getElementById("todo-input");

    if (input.value === "") {
        return;
    }

    const id = getLastID();
    createTodoUI(id, todosElement, input.value, false);

    sendTodos(input);

    input.value = "";
}

const todosElement = createTodosElement();
createBasicUI(todosElement);
getTodos(todosElement);
