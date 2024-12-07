document.addEventListener('DOMContentLoaded', () => {
    const app = document.getElementById('app');
    
    // タスク入力フォームの作成
    const form = document.createElement('form');
    form.innerHTML = `
        <div style="margin-bottom: 20px">
            <input type="text" id="taskInput" placeholder="新しいタスクを入力" style="padding: 8px; width: 70%">
            <button type="submit" style="padding: 8px; margin-left: 10px">追加</button>
        </div>
    `;
    
    // タスクリストの表示領域
    const taskList = document.createElement('div');
    app.appendChild(form);
    app.appendChild(taskList);
    
    // タスクの読み込みと表示
    const loadTasks = async () => {
        const response = await fetch('/api/tasks');
        const tasks = await response.json();
        
        taskList.innerHTML = tasks.map(task => `
            <div style="padding: 10px; margin: 5px 0; background: #f9f9f9; border-radius: 4px">
                <input type="checkbox" ${task.done ? 'checked' : ''}>
                <span>${task.title}</span>
                <small style="color: #666">作成日: ${new Date(task.created_at).toLocaleString()}</small>
            </div>
        `).join('');
    };
    
    // 新しいタスクの追加
    form.onsubmit = async (e) => {
        e.preventDefault();
        const input = document.getElementById('taskInput');
        const title = input.value.trim();
        
        if (title) {
            await fetch('/api/tasks', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ title, done: false }),
            });
            
            input.value = '';
            loadTasks();
        }
    };
    
    // 初期表示
    loadTasks();
});