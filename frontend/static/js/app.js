// 文件上传功能
document.addEventListener('DOMContentLoaded', () => {
    const uploadArea = document.getElementById('uploadArea');
    const fileInput = document.getElementById('fileInput');
    const fileList = document.getElementById('fileList');
    
    // 拖拽上传功能
    uploadArea.addEventListener('dragover', (e) => {
        e.preventDefault();
        uploadArea.style.borderColor = '#1976D2';
    });

    uploadArea.addEventListener('dragleave', () => {
        uploadArea.style.borderColor = '#ccc';
    });

    uploadArea.addEventListener('drop', (e) => {
        e.preventDefault();
        uploadArea.style.borderColor = '#ccc';
        if (e.dataTransfer.files.length) {
            handleFiles(e.dataTransfer.files);
        }
    });

    // 点击上传区域触发文件选择
    uploadArea.addEventListener('click', () => {
        fileInput.click();
    });

    // 文件选择变化处理
    fileInput.addEventListener('change', () => {
        if (fileInput.files && fileInput.files.length > 0) {
            handleFiles(fileInput.files);
            fileInput.value = ''; // 重置选择，允许重复选择相同文件
        }
    });

    // 处理文件上传
    async function handleFiles(files) {
        const formData = new FormData();
        for (let i = 0; i < files.length; i++) {
            formData.append('files', files[i]);
        }

        try {
            const response = await fetch('/api/upload', {
                method: 'POST',
                body: formData
            });
            
            const result = await response.json();
            if (response.ok) {
                alert('上传成功');
                refreshFileList();
            } else {
                alert(`上传失败: ${result.message}`);
            }
        } catch (error) {
            console.error('上传错误:', error);
            alert('上传过程中发生错误');
        }
    }

    // 刷新文件列表
    async function refreshFileList() {
        const loading = document.getElementById('loading');
        try {
            loading.style.display = 'block';
            fileList.innerHTML = '';
            
            const response = await fetch('/api/files');
            if (!response.ok) {
                throw new Error('获取文件列表失败');
            }
            const files = await response.json();
            console.log('文件列表API响应:', files); // 调试日志
            
            loading.style.display = 'none';
            files.forEach(file => {
                console.log('单个文件数据:', file); // 调试日志
                const li = document.createElement('li');
                li.className = 'file-item';
                li.innerHTML = `
                    <span>${file.Name}</span>
                    <span>
                        ${new Date(file.UpdatedAt).toLocaleString()} | 
                        ${formatFileSize(file.Size)}
                    </span>
                    <span>
                        <button data-id="${file.ID}" class="download-btn">下载</button>
                        <button data-id="${file.ID}" class="delete-btn">删除</button>
                    </span>
                `;
                fileList.appendChild(li);
            });

            // 使用事件委托处理按钮点击
            fileList.addEventListener('click', (e) => {
                const btn = e.target.closest('button');
                if (!btn) return;
                
                const fileId = btn.getAttribute('data-id');
                if (!fileId) {
                    console.error('无法获取文件ID', btn);
                    return;
                }

                if (btn.classList.contains('download-btn')) {
                    downloadFile(fileId);
                } else if (btn.classList.contains('delete-btn')) {
                    deleteFile(fileId);
                }
            });
        } catch (error) {
            console.error('获取文件列表失败:', error);
        }
    }

    // 初始化文件列表
    refreshFileList();
});

// 辅助函数
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// 全局函数声明
const refreshFileList = async () => {
    const loading = document.getElementById('loading');
    try {
        loading.style.display = 'block';
        fileList.innerHTML = '';
        
        const response = await fetch('/api/files');
        if (!response.ok) {
            throw new Error('获取文件列表失败');
        }
        const files = await response.json();
        console.log('文件列表API响应:', files);
        
        loading.style.display = 'none';
        files.forEach(file => {
            console.log('单个文件数据:', file);
            const li = document.createElement('li');
            li.className = 'file-item';
            li.innerHTML = `
                <span>${file.Name}</span>
                <span>
                    ${new Date(file.UpdatedAt).toLocaleString()} | 
                    ${formatFileSize(file.Size)}
                </span>
                <span>
                    <button data-id="${file.ID}" class="download-btn">下载</button>
                    <button data-id="${file.ID}" class="delete-btn">删除</button>
                </span>
            `;
            fileList.appendChild(li);
        });
    } catch (error) {
        console.error('获取文件列表失败:', error);
    }
};

// 文件操作函数
async function downloadFile(fileId) {
    if (!fileId || fileId === 'undefined') {
        alert('无效的文件ID');
        return;
    }
    // 直接使用后端下载URL
    window.location.href = `/api/download/${fileId}`;
}

async function deleteFile(fileId) {
    if (confirm('确定要删除此文件吗？')) {
        try {
            const response = await fetch(`/api/files/${fileId}`, {
                method: 'DELETE'
            });
            
            if (response.ok) {
                alert('删除成功');
                await refreshFileList();
            } else {
                const result = await response.json();
            }
        } catch (error) {
            console.error('删除错误:', error);
            alert('删除过程中发生错误');
        }
    }
}

// 初始化
document.addEventListener('DOMContentLoaded', () => {
    // 原有初始化代码...
    refreshFileList();
});
