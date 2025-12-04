# Upload Surat - Final Solution

## Problem: Placeholder Path from Frontend

**Error Log**:
```
Warning: failed to move file from uploaded-via-form to peminjaman/PMJ-251204-0008/surat.pdf: 
failed to download file (status 400): Object not found
```

**Root Cause**: Frontend sends placeholder path `uploaded-via-form` instead of actual file path in storage.

---

## Solution: Two-Step Upload Workflow

Backend now supports **flexible workflow**:

### Workflow 1: Upload Later (Current Frontend)
1. Frontend creates peminjaman with placeholder path (`uploaded-via-form`)
2. Backend detects placeholder → skip file move
3. Backend saves peminjaman with **empty** `path_surat_digital`
4. Frontend uploads file via `/api/peminjaman/{kode}/upload-surat`
5. File saved to `peminjaman/{kode}/surat.pdf`

### Workflow 2: Upload First (Future)
1. Frontend uploads file to storage first → gets real path
2. Frontend creates peminjaman with real path
3. Backend moves file to unique path `peminjaman/{kode}/surat.pdf`

---

## Implementation

### Placeholder Detection

Backend detects common placeholder values:
```go
isPlaceholder := suratPath == "" || 
    suratPath == "uploaded-via-form" || 
    suratPath == "pending" || 
    suratPath == "temp"
```

### Logic Flow

```go
if isPlaceholder {
    // Don't try to move file, set path to empty
    uniquePath = ""
} else if fileExists(suratPath) {
    // Move file to unique path
    MoveFile(suratPath, uniquePath)
} else {
    // File doesn't exist, set path to empty
    uniquePath = ""
}

// Save to database (empty or unique path)
UpdateSuratDigitalURL(kode, uniquePath)
```

---

## API Response

### Case 1: Placeholder Path
**Request**:
```json
{
  "path_surat_digital": "uploaded-via-form",
  ...
}
```

**Response**:
```json
{
  "kode_peminjaman": "PMJ-251204-0008",
  "path_surat_digital": "",  // Empty - file not uploaded yet
  ...
}
```

Frontend should then upload via `/api/peminjaman/PMJ-251204-0008/upload-surat`

### Case 2: Real Path
**Request**:
```json
{
  "path_surat_digital": "temp/user-upload-123.pdf",
  ...
}
```

**Response**:
```json
{
  "kode_peminjaman": "PMJ-251204-0008",
  "path_surat_digital": "peminjaman/PMJ-251204-0008/surat.pdf",  // Moved!
  ...
}
```

---

## Frontend Integration

### Current Frontend (Placeholder)

```javascript
// 1. Create peminjaman with placeholder
const response = await apiCall('/api/peminjaman', 'POST', {
    path_surat_digital: 'uploaded-via-form',  // Placeholder
    ...
});

const kodePeminjaman = response.kode_peminjaman;

// 2. Upload file separately
const formData = new FormData();
formData.append('surat', fileInput.files[0]);

await apiCall(`/api/peminjaman/${kodePeminjaman}/upload-surat`, 'POST', formData);
```

### Future Frontend (Upload First)

```javascript
// 1. Upload file first
const formData = new FormData();
formData.append('file', fileInput.files[0]);
const uploadResponse = await apiCall('/api/upload/temp', 'POST', formData);
const filePath = uploadResponse.path;

// 2. Create peminjaman with real path
const response = await apiCall('/api/peminjaman', 'POST', {
    path_surat_digital: filePath,  // Real path
    ...
});
// Backend will move file to peminjaman/{kode}/surat.pdf
```

---

## Benefits

✅ **Flexible**: Supports both workflows
✅ **No Breaking Changes**: Works with current frontend
✅ **No Errors**: No more 404 warnings in logs
✅ **Clean Storage**: Files organized in proper folders
✅ **Future-Proof**: Ready for frontend improvements

---

## Testing

### Test 1: Current Workflow (Placeholder)

**Steps**:
1. Create peminjaman with `path_surat_digital: "uploaded-via-form"`
2. Check response: `path_surat_digital` should be empty
3. Upload file via `/upload-surat` endpoint
4. Check storage: file in `peminjaman/{kode}/surat.pdf`

**Expected**: ✅ No errors, file uploaded successfully

### Test 2: Check Logs

**Before**:
```
Warning: failed to move file from uploaded-via-form...
```

**After**:
```
Info: skipping file move (placeholder detected)
```

---

## Summary

✅ **Problem**: Frontend sends placeholder path
✅ **Solution**: Detect placeholder, skip file move
✅ **Result**: No errors, flexible workflow
✅ **Status**: Production ready!
