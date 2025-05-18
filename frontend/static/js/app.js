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
                refreshFileList('after upload');
            } else {
                alert(`上传失败: ${result.message}`);
            }
        } catch (error) {
            console.error('上传错误:', error);
            alert('上传过程中发生错误');
        }
    }

    // 刷新文件列表
    async function refreshFileList(caller = 'unknown') {
        console.log(`refreshFileList called from: ${caller}`);
        const loading = document.getElementById('loading');
        const fileList = document.getElementById('fileList');
        
        // 清除现有事件监听
        const newFileList = fileList.cloneNode(false);
        fileList.parentNode.replaceChild(newFileList, fileList);
        
        try {
            loading.style.display = 'block';
            newFileList.innerHTML = '';
            
            const response = await fetch('/api/files?t=' + Date.now()); // 防止缓存
            if (!response.ok) {
                throw new Error(`获取文件列表失败 (状态码: ${response.status})`);
            }
            const files = await response.json();
            
            if (!Array.isArray(files)) {
                throw new Error('服务器返回无效的文件列表格式');
            }

            // 去重处理
            const uniqueFiles = [];
            const seenIds = new Set();
            files.forEach(file => {
                if (!seenIds.has(file.ID)) {
                    seenIds.add(file.ID);
                    uniqueFiles.push(file);
                }
            });

            loading.style.display = 'none';
            
            if (uniqueFiles.length === 0) {
                const li = document.createElement('li');
                li.className = 'file-item empty';
                li.textContent = '暂无文件';
                newFileList.appendChild(li);
                return;
            }

            uniqueFiles.forEach(file => {
                const li = document.createElement('li');
                li.className = 'file-item';
                li.innerHTML = `
                    <span>${file.Name || '未命名文件'}</span>
                    <span>
                        ${file.UpdatedAt ? new Date(file.UpdatedAt).toLocaleString() : '未知时间'} | 
                        ${file.Size ? formatFileSize(file.Size) : '未知大小'}
                    </span>
                    <span>
                        <button data-id="${file.ID}" class="download-btn">下载</button>
                        <button data-id="${file.ID}" class="delete-btn">删除</button>
                    </span>
                `;
                newFileList.appendChild(li);
            });

            // 使用事件委托处理按钮点击
            newFileList.addEventListener('click', (e) => {
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
            loading.style.display = 'none';
            const li = document.createElement('li');
            li.className = 'file-item error';
            li.textContent = `加载失败: ${error.message}`;
            fileList.appendChild(li);
        }
    }

    // 初始化文件列表
    refreshFileList('initial load');
});

// 辅助函数
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}



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
                await refreshFileList('after delete');
            } else {
                const result = await response.json();
            }
        } catch (error) {
            console.error('删除错误:', error);
            alert('删除过程中发生错误');
        }
    }
}


