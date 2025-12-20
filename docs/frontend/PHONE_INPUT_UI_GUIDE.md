# Frontend Guide: Phone Number Input - Best Practices

## 🎨 Recommended UI/UX for Indonesian Phone Numbers

### 📱 Best Practice Implementation

```html
<!-- Enhanced Phone Number Input -->
<div class="form-group">
    <label for="no_hp">
        Nomor WhatsApp 
        <span class="optional">(Opsional)</span>
    </label>
    <input 
        type="tel" 
        id="no_hp" 
        name="no_hp"
        placeholder="+62 812 3456 7890"
        pattern="^(\+62|62|0)[0-9]{9,13}$"
        title="Format: +62xxxxxxxxxx atau 08xxxxxxxxxx"
    >
    <small class="hint">Format: +62 atau 08 diikuti 9-13 digit</small>
</div>
```

### ✨ Enhanced Version with Auto-formatting

```html
<div class="form-group">
    <label for="no_hp">
        <i>📱</i> Nomor WhatsApp 
        <span class="badge-optional">Opsional</span>
    </label>
    <div class="input-wrapper">
        <span class="prefix">+62</span>
        <input 
            type="tel" 
            id="no_hp" 
            name="no_hp"
            placeholder="812 3456 7890"
            maxlength="15"
            oninput="formatPhoneNumber(this)"
        >
    </div>
    <div class="hint">
        <small>💬 Untuk notifikasi WhatsApp</small>
    </div>
</div>
```

### 🎨 CSS Styling

```css
/* Phone number input styling */
.form-group {
    margin-bottom: 20px;
}

.form-group label {
    display: block;
    margin-bottom: 8px;
    font-weight: 600;
    color: #333;
}

.badge-optional {
    display: inline-block;
    padding: 2px 8px;
    font-size: 11px;
    font-weight: 500;
    color: #6c757d;
    background: #f8f9fa;
    border: 1px solid #dee2e6;
    border-radius: 12px;
    margin-left: 8px;
}

.input-wrapper {
    display: flex;
    align-items: center;
    border: 2px solid #dee2e6;
    border-radius: 8px;
    overflow: hidden;
    transition: all 0.3s ease;
}

.input-wrapper:focus-within {
    border-color: #FF6B35;
    box-shadow: 0 0 0 3px rgba(255, 107, 53, 0.1);
}

.input-wrapper .prefix {
    padding: 12px 16px;
    background: #f8f9fa;
    font-weight: 600;
    color: #495057;
    border-right: 2px solid #dee2e6;
    user-select: none;
}

.input-wrapper input {
    flex: 1;
    border: none;
    padding: 12px 16px;
    font-size: 16px;
    outline: none;
}

.hint {
    margin-top: 6px;
    color: #6c757d;
    font-size: 13px;
}

.hint i {
    margin-right: 4px;
}
```

### 🔧 JavaScript - Auto Formatting

```javascript
// Format phone number as user types
function formatPhoneNumber(input) {
    // Remove all non-digit characters
    let value = input.value.replace(/\D/g, '');
    
    // Limit to 13 digits (for +62 format)
    value = value.substring(0, 13);
    
    // Add spaces for readability
    if (value.length > 3) {
        value = value.substring(0, 3) + ' ' + value.substring(3);
    }
    if (value.length > 7) {
        value = value.substring(0, 7) + ' ' + value.substring(7);
    }
    if (value.length > 12) {
        value = value.substring(0, 12) + ' ' + value.substring(12);
    }
    
    input.value = value;
}

// Convert to backend format before submit
function preparePhoneNumber(phoneInput) {
    if (!phoneInput) return null;
    
    // Remove spaces and dashes
    let phone = phoneInput.replace(/[\s-]/g, '');
    
    // If starts with 0, replace with +62
    if (phone.startsWith('0')) {
        phone = '+62' + phone.substring(1);
    }
    // If starts with 62, add +
    else if (phone.startsWith('62')) {
        phone = '+' + phone;
    }
    // If doesn't start with +, assume +62
    else if (!phone.startsWith('+')) {
        phone = '+62' + phone;
    }
    
    return phone;
}

// On form submit
document.getElementById('registerForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const formData = new FormData(e.target);
    const phoneInput = formData.get('no_hp');
    
    const registerData = {
        nama: formData.get('nama'),
        email: formData.get('email'),
        password: formData.get('password'),
        role: 'MAHASISWA',
        no_hp: phoneInput ? preparePhoneNumber(phoneInput) : null,
        organisasi_kode: formData.get('organisasi_kode') || null
    };
    
    // Submit to backend
    try {
        const response = await fetch(`${API_BASE_URL}/auth/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(registerData)
        });
        
        if (!response.ok) throw new Error('Registration failed');
        
        alert('✅ Registration successful!');
        window.location.href = '/login.html';
    } catch (error) {
        alert('❌ ' + error.message);
    }
});
```

---

## 🎯 Three Implementation Options

### Option 1: Simple (Use Suggestion from Friend) ⭐ RECOMMENDED

```html
<input 
    type="tel" 
    name="no_hp"
    placeholder="+62"
    pattern="^(\+62|0)[0-9]{9,13}$"
>
```

**Pros:**
- ✅ Simple & clean
- ✅ Clear indication of Indonesian format
- ✅ No JavaScript needed
- ✅ Fast implementation

**Cons:**
- ⚠️ User needs to type full number
- ⚠️ No auto-formatting

---

### Option 2: With Prefix Style (Medium)

```html
<div class="phone-input">
    <span class="prefix">+62</span>
    <input 
        type="tel" 
        name="no_hp"
        placeholder="812 3456 7890"
        oninput="formatPhoneNumber(this)"
    >
</div>
```

**Pros:**
- ✅ Visual prefix (+62 always visible)
- ✅ User only types numbers
- ✅ Professional look
- ✅ Auto-formatting (with JS)

**Cons:**
- ⚠️ Needs CSS styling
- ⚠️ Needs JavaScript for formatting

---

### Option 3: Advanced with Validation (Complex)

```html
<!-- With country code dropdown -->
<div class="phone-group">
    <select class="country-code">
        <option value="+62" selected>🇮🇩 +62</option>
        <option value="+65">🇸🇬 +65</option>
        <option value="+60">🇲🇾 +60</option>
    </select>
    <input 
        type="tel" 
        name="no_hp"
        placeholder="812 3456 7890"
    >
</div>
```

**Pros:**
- ✅ International support
- ✅ Flexible for different countries
- ✅ Premium look

**Cons:**
- ⚠️ Overkill for Indonesia-only app
- ⚠️ More complex code

---

## 💡 **RECOMMENDATION: Use Option 1**

Sesuai saran teman kamu, **paling simple dan effective**:

```html
<label>
    Nomor WhatsApp <span class="optional">(Opsional)</span>
</label>
<input 
    type="tel" 
    name="no_hp"
    placeholder="+62 812 3456 7890"
    pattern="^(\+62|62|0)[0-9]{9,13}$"
    title="Format: +62812xxxxxxx atau 08xxxxxxxx"
>
<small>💬 Untuk notifikasi WhatsApp</small>
```

**Why?**
- ✅ User langsung paham format (+62)
- ✅ No complex JavaScript
- ✅ Clear & simple UX
- ✅ Works on all devices

---

## 📝 Complete Example - Registration Form

```html
<!DOCTYPE html>
<html>
<head>
    <title>Register - Sarpras</title>
    <style>
        .form-container {
            max-width: 400px;
            margin: 50px auto;
            padding: 30px;
            background: white;
            border-radius: 12px;
            box-shadow: 0 4px 20px rgba(0,0,0,0.1);
        }
        
        .form-group {
            margin-bottom: 20px;
        }
        
        label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
            color: #333;
        }
        
        .optional {
            font-weight: 400;
            color: #6c757d;
            font-size: 14px;
        }
        
        input, select {
            width: 100%;
            padding: 12px 16px;
            border: 2px solid #dee2e6;
            border-radius: 8px;
            font-size: 16px;
            transition: all 0.3s;
        }
        
        input:focus, select:focus {
            outline: none;
            border-color: #FF6B35;
            box-shadow: 0 0 0 3px rgba(255,107,53,0.1);
        }
        
        input[type="tel"]::placeholder {
            color: #adb5bd;
            font-style: italic;
        }
        
        small {
            display: block;
            margin-top: 6px;
            color: #6c757d;
            font-size: 13px;
        }
        
        button {
            width: 100%;
            padding: 14px;
            background: linear-gradient(135deg, #FF6B35 0%, #F7931E 100%);
            color: white;
            border: none;
            border-radius: 8px;
            font-weight: 600;
            font-size: 16px;
            cursor: pointer;
            transition: transform 0.2s;
        }
        
        button:hover {
            transform: translateY(-2px);
        }
    </style>
</head>
<body>
    <div class="form-container">
        <h2>📝 Registrasi</h2>
        
        <form id="registerForm">
            <div class="form-group">
                <label>Nama Lengkap</label>
                <input type="text" name="nama" required>
            </div>
            
            <div class="form-group">
                <label>Email</label>
                <input type="email" name="email" required>
            </div>
            
            <div class="form-group">
                <label>Password</label>
                <input type="password" name="password" required minlength="6">
            </div>
            
            <!-- PHONE NUMBER INPUT (As Suggested) -->
            <div class="form-group">
                <label>
                    Nomor WhatsApp 
                    <span class="optional">(Opsional)</span>
                </label>
                <input 
                    type="tel" 
                    name="no_hp"
                    placeholder="+62 812 3456 7890"
                    pattern="^(\+62|62|0)[0-9]{9,13}$"
                >
                <small>💬 Untuk menerima notifikasi WhatsApp</small>
            </div>
            
            <div class="form-group">
                <label>Organisasi <span class="optional">(Opsional)</span></label>
                <select name="organisasi_kode">
                    <option value="">Pilih Organisasi</option>
                    <!-- Options dari API -->
                </select>
            </div>
            
            <button type="submit">Daftar</button>
        </form>
    </div>
    
    <script>
        const API_BASE_URL = 'http://localhost:8000/api';
        
        document.getElementById('registerForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const formData = new FormData(e.target);
            
            const registerData = {
                nama: formData.get('nama'),
                email: formData.get('email'),
                password: formData.get('password'),
                role: 'MAHASISWA',
                no_hp: formData.get('no_hp') || null,
                organisasi_kode: formData.get('organisasi_kode') || null
            };
            
            try {
                const response = await fetch(`${API_BASE_URL}/auth/register`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(registerData)
                });
                
                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error);
                }
                
                alert('✅ Registrasi berhasil!');
                window.location.href = '/login.html';
            } catch (error) {
                alert('❌ Registrasi gagal: ' + error.message);
            }
        });
    </script>
</body>
</html>
```

---

## ✅ Summary

Sesuai saran teman kamu:

```html
<input 
    type="tel" 
    name="no_hp"
    placeholder="+62 812 3456 7890"
>
```

**Perfect karena:**
- ✅ Simple & jelas
- ✅ User paham format Indonesia
- ✅ No kompleksitas berlebih
- ✅ Mobile-friendly
- ✅ Accessibility-friendly

---

*Best Practice for Indonesian Phone Number Input ✨*
