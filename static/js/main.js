// グローバルなtoken変数
let token = localStorage.getItem('token');

// 画面表示の制御
function showWelcome() {
    document.getElementById('welcome-section').style.display = 'block';
    document.getElementById('login-section').style.display = 'none';
    document.getElementById('register-section').style.display = 'none';
    document.getElementById('task-section').style.display = 'none';
}

function showLoginForm() {
    document.getElementById('welcome-section').style.display = 'none';
    document.getElementById('login-section').style.display = 'block';
    document.getElementById('register-section').style.display = 'none';
}

function showRegisterForm() {
    document.getElementById('welcome-section').style.display = 'none';
    document.getElementById('login-section').style.display = 'none';
    document.getElementById('register-section').style.display = 'block';
}

window.showTaskSection = function() {
    document.getElementById('welcome-section').style.display = 'none';
    document.getElementById('login-section').style.display = 'none';
    document.getElementById('register-section').style.display = 'none';
    document.getElementById('task-section').style.display = 'block';
}

// 認証関連の関数
window.register = async function() {
    const username = document.getElementById('register-username').value;
    const password = document.getElementById('register-password').value;
    const confirmPassword = document.getElementById('register-password-confirm').value;

    if (password !== confirmPassword) {
        alert('パスワードが一致しません');
        return;
    }

    try {
        const response = await fetch('/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password }),
        });

        const data = await response.json();
        if (response.ok) {
            alert('登録が完了しました。ログインしてください。');
            showLoginForm();
        } else {
            alert(data.error);
        }
    } catch (error) {
        alert('エラーが発生しました');
    }
}

window.login = async function() {
    const username = document.getElementById('login-username').value;
    const password = document.getElementById('login-password').value;

    try {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password }),
        });

        const data = await response.json();
        if (response.ok) {
            token = data.token;
            localStorage.setItem('token', token);
            showTaskSection();
            loadTasks();
        } else {
            alert(data.error);
        }
    } catch (error) {
        alert('エラーが発生しました');
    }
}

// タスク関連の関数
window.createTask = async function() {
    const title = document.getElementById('task-title').value;
    const description = document.getElementById('task-description').value;
    const dueDate = document.getElementById('task-due-date').value;
    const priority = document.getElementById('task-priority').value;

    try {
        const response = await fetch('/api/tasks', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': token,
            },
            body: JSON.stringify({
                title,
                description,
                due_date: new Date(dueDate),
                priority: parseInt(priority),
                done: false,
            }),
        });

        if (response.ok) {
            clearTaskForm();
            loadTasks();
        } else {
            const data = await response.json();
            alert(data.error);
        }
    } catch (error) {
        alert('エラーが発生しました');
    }
}

window.loadTasks = async function() {
    try {
        const response = await fetch('/api/tasks', {
            headers: {
                'Authorization': token,
            },
        });

        const tasks = await response.json();
        const taskList = document.getElementById('task-list');
        taskList.innerHTML = tasks.map(task => `
            <div class="task-item ${getPriorityClass(task.priority)}">
                <h3>${task.title}</h3>
                <p>${task.description}</p>
                <p>期限: ${new Date(task.due_date).toLocaleDateString()}</p>
                <div class="task-controls">
                    <button onclick="toggleTaskStatus(${task.ID})" class="complete-btn">
                        ${task.done ? '未完了に戻す' : '完了にする'}
                    </button>
                    <button onclick="deleteTask(${task.ID})" class="delete-btn">削除</button>
                </div>
            </div>
        `).join('');
    } catch (error) {
        alert('タスクの読み込みに失敗しました');
    }
}

window.toggleTaskStatus = async function(taskId) {
    try {
        const response = await fetch(`/api/tasks/${taskId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': token,
            },
            body: JSON.stringify({ done: true }),
        });

        if (response.ok) {
            loadTasks();
        }
    } catch (error) {
        alert('エラーが発生しました');
    }
}

window.deleteTask = async function(taskId) {
    if (!confirm('本当に削除しますか？')) return;

    try {
        const response = await fetch(`/api/tasks/${taskId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': token,
            },
        });

        if (response.ok) {
            loadTasks();
        }
    } catch (error) {
        alert('エラーが発生しました');
    }
}

// ユーティリティ関数
function clearTaskForm() {
    document.getElementById('task-title').value = '';
    document.getElementById('task-description').value = '';
    document.getElementById('task-due-date').value = '';
    document.getElementById('task-priority').value = '1';
}

function getPriorityClass(priority) {
    switch (priority) {
        case 3: return 'high-priority';
        case 2: return 'medium-priority';
        default: return '';
    }
}

window.logout = function() {
    localStorage.removeItem('token');
    showWelcome();
}

// 初期化
if (token) {
    showTaskSection();
    loadTasks();
} else {
    showWelcome();
}