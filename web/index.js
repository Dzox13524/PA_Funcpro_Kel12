const API_BASE = "http://localhost:8080/api/v1";
let authToken = localStorage.getItem("authToken") || null;
let currentUser = null;
let searchTimeout;

document.addEventListener("DOMContentLoaded", () => {
  if (authToken) {
    loadUserProfile();
  }
  loadProducts();
  loadQuestions();
  loadAlerts();
  createParticles();
  loadProductMeta();
});

async function buyProduct(id) {
  if (!authToken) {
    showToast("Silakan login untuk membeli produk", "error");
    openLoginModal();
    return;
  }

  try {
    const res = await fetch(`${API_BASE}/market/products/${id}`);
    const data = await res.json();

    if (res.ok && data.data) {
      const p = data.data;

      if (p.stock < 1) {
        showToast("Maaf, stok produk ini habis!", "error");
        return;
      }

      document.getElementById("buyProductId").value = p.id;
      document.getElementById("buyProductName").textContent = p.name;
      document.getElementById(
        "buyProductPriceDisplay"
      ).textContent = `Rp ${Number(p.price).toLocaleString()}/unit`;
      document.getElementById("buyProductStock").textContent = p.stock;
      document.getElementById("buyProductPriceUnit").value = p.price;

      document.getElementById("buyQuantity").value = 1;
      document.getElementById("buyQuantity").max = p.stock;
      document.getElementById("buyNote").value = "";

      calculateTotal();

      closeModal("productDetailModal");
      document.getElementById("buyModal").classList.add("active");
    }
  } catch (err) {
    showToast("Gagal memuat data produk", "error");
  }
}

// ==========================================
// LOGIKA RESERVASI
// ==========================================

// 1. Buka Modal Reservasi
async function openReservationModal(id) {
    if (!authToken) {
        showToast("Silakan login untuk melakukan reservasi", "error");
        openLoginModal();
        return;
    }

    try {
        const res = await fetch(`${API_BASE}/market/products/${id}`);
        const data = await res.json();

        if (res.ok && data.data) {
            const p = data.data;

            // Isi Data ke Modal
            document.getElementById("resProductId").value = p.id;
            document.getElementById("resProductName").textContent = p.name;
            document.getElementById("resProductStock").textContent = p.stock;
            
            // Reset Form
            document.getElementById("resQuantity").value = 1;
            // Kita batasi input max sesuai stok (opsional, tergantung kebijakan)
            document.getElementById("resQuantity").max = p.stock; 
            document.getElementById("resNote").value = "";

            // Tampilkan Modal
            closeModal("productDetailModal");
            document.getElementById("reservationModal").classList.add("active");
        }
    } catch (err) {
        showToast("Gagal memuat data produk", "error");
    }
}

function adjustResQty(change) {
    const qtyInput = document.getElementById("resQuantity");
    const stockMax = parseInt(document.getElementById("resProductStock").textContent);
    
    let currentQty = parseInt(qtyInput.value) || 0;
    let newQty = currentQty + change;

    if (newQty < 1) newQty = 1;
    if (newQty > stockMax) {
        newQty = stockMax;
        showToast(`Maksimal reservasi ${stockMax} unit`, "error");
    }

    qtyInput.value = newQty;
}

async function handleReservationSubmit(e) {
    e.preventDefault();
    
    const btn = document.getElementById("btnConfirmRes");
    const productId = document.getElementById("resProductId").value;
    const quantity = parseInt(document.getElementById("resQuantity").value);
    const note = document.getElementById("resNote").value;

    if (!note) {
        showToast("Mohon isi catatan reservasi", "error");
        return;
    }

    btn.innerHTML = '<div class="loading" style="width:15px; height:15px; border-color:black;"></div> Memproses...';
    btn.disabled = true;

    try {
        const res = await fetch(`${API_BASE}/market/reservations`, {
            method: "POST",
            headers: {
                Authorization: `Bearer ${authToken}`,
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                product_id: productId,
                quantity: quantity,
                note: note
            }),
        });

        const data = await res.json();

        if (res.ok) {
            showToast("Reservasi berhasil diajukan! 📅");
            closeModal("reservationModal");
            
            if (typeof openProfile === "function" && typeof switchProfileTab === "function") {
                openProfile();
                switchProfileTab('my-reservations');
            }
        } else {
            showToast(data.message || "Gagal membuat reservasi", "error");
        }
    } catch (err) {
        console.error(err);
        showToast("Terjadi kesalahan koneksi", "error");
    } finally {
        btn.innerHTML = "Konfirmasi Reservasi";
        btn.disabled = false;
    }
}

function calculateTotal() {
  const qtyInput = document.getElementById("buyQuantity");
  const priceUnit = parseInt(
    document.getElementById("buyProductPriceUnit").value
  );
  const stockMax = parseInt(
    document.getElementById("buyProductStock").textContent
  );

  let qty = parseInt(qtyInput.value);

  if (qty < 1) qty = 1;
  if (qty > stockMax) {
    qty = stockMax;
    showToast(`Maksimal pembelian ${stockMax} unit`, "error");
  }

  qtyInput.value = qty;

  const total = qty * priceUnit;
  document.getElementById("buyTotalPrice").textContent = `Rp ${Number(
    total
  ).toLocaleString()}`;
}

function adjustQty(change) {
  const qtyInput = document.getElementById("buyQuantity");
  let currentQty = parseInt(qtyInput.value) || 0;
  qtyInput.value = currentQty + change;
  calculateTotal();
}

async function handleOrderSubmit(e) {
  e.preventDefault();

  const btn = document.getElementById("btnConfirmOrder");
  const productId = document.getElementById("buyProductId").value;
  const quantity = parseInt(document.getElementById("buyQuantity").value);
  const note = document.getElementById("buyNote").value;

  btn.innerHTML =
    '<div class="loading" style="width:15px; height:15px;"></div> Memproses...';
  btn.disabled = true;

  try {
    const res = await fetch(`${API_BASE}/market/orders`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${authToken}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        product_id: productId,
        quantity: quantity,
        note: note,
      }),
    });

    const data = await res.json();

    if (res.ok) {
      showToast("Pesanan berhasil dibuat! 📦");
      closeModal("buyModal");

      if (typeof openProfile === "function") {
        openProfile();
        switchProfileTab("my-orders");
      } else {
        loadProducts();
      }
    } else {
      showToast(data.message || "Gagal membuat pesanan", "error");
    }
  } catch (err) {
    console.error(err);
    showToast("Terjadi kesalahan koneksi", "error");
  } finally {
    btn.innerHTML = "Buat Pesanan";
    btn.disabled = false;
  }
}

async function loadProductMeta() {
  const wilayahJatim = [
    "Bangkalan",
    "Banyuwangi",
    "Blitar",
    "Bojonegoro",
    "Bondowoso",
    "Gresik",
    "Jember",
    "Jombang",
    "Kediri",
    "Lamongan",
    "Lumajang",
    "Madiun",
    "Magetan",
    "Malang",
    "Mojokerto",
    "Nganjuk",
    "Ngawi",
    "Pacitan",
    "Pamekasan",
    "Pasuruan",
    "Ponorogo",
    "Probolinggo",
    "Sampang",
    "Sidoarjo",
    "Situbondo",
    "Sumenep",
    "Trenggalek",
    "Tuban",
    "Tulungagung",
    "Kota Batu",
    "Kota Blitar",
    "Kota Kediri",
    "Kota Madiun",
    "Kota Malang",
    "Kota Mojokerto",
    "Kota Pasuruan",
    "Kota Probolinggo",
    "Kota Surabaya",
  ];

  const jenisTanaman = [
    "Padi",
    "Jagung",
    "Kedelai",
    "Cabai Rawit",
    "Cabai Merah",
    "Bawang Merah",
    "Bawang Putih",
    "Tomat",
    "Wortel",
    "Kubis",
    "Kentang",
    "Tebu",
    "Tembakau",
    "Kopi",
    "Kakao",
  ];

  try {
    const cropSelect = document.getElementById("productCrop");
    const regionSelect = document.getElementById("productRegion");
    const searchCropSelect = document.getElementById("searchCrop");
    const searchRegionSelect = document.getElementById("searchRegion");
    const editCropSelect = document.getElementById("editProductCrop");
    const editRegionSelect = document.getElementById("editProductRegion");
    const pestCitySelect = document.getElementById("pestCity");

    const cropOpts = jenisTanaman
      .map((item) => `<option value="${item}">${item}</option>`)
      .join("");
    const regionOpts = wilayahJatim
      .map((item) => `<option value="${item}">${item}</option>`)
      .join("");

    if (cropSelect)
      cropSelect.innerHTML =
        `<option value="">Pilih Jenis Tanaman</option>` + cropOpts;
    if (regionSelect)
      regionSelect.innerHTML =
        `<option value="">Pilih Wilayah</option>` + regionOpts;
    if (searchCropSelect)
      searchCropSelect.innerHTML =
        `<option value="">Semua Tanaman</option>` + cropOpts;
    if (searchRegionSelect)
      searchRegionSelect.innerHTML =
        `<option value="">Semua Wilayah</option>` + regionOpts;
    if (editCropSelect)
      editCropSelect.innerHTML =
        `<option value="">Pilih Jenis Tanaman</option>` + cropOpts;
    if (editRegionSelect)
      editRegionSelect.innerHTML =
        `<option value="">Pilih Wilayah</option>` + regionOpts;
    if (pestCitySelect) {
      pestCitySelect.innerHTML =
        `<option value="">Pilih Wilayah</option>` + regionOpts;
    }
  } catch (err) {
    console.error("Gagal memuat data meta:", err);
  }
}

window.addEventListener("scroll", () => {
  const navbar = document.getElementById("navbar");
  if (window.scrollY > 50) {
    navbar.classList.add("scrolled");
  } else {
    navbar.classList.remove("scrolled");
  }
});

function toggleMenu() {
  document.getElementById("navMenu").classList.toggle("active");
}

function scrollToSection(id) {
  document.getElementById(id).scrollIntoView({ behavior: "smooth" });
}

function createParticles() {
  const container = document.getElementById("particles");
  for (let i = 0; i < 30; i++) {
    const particle = document.createElement("div");
    particle.className = "particle";
    particle.style.left = Math.random() * 100 + "%";
    particle.style.top = Math.random() * 100 + "%";
    particle.style.animationDelay = Math.random() * 20 + "s";
    particle.style.animationDuration = Math.random() * 10 + 15 + "s";
    container.appendChild(particle);
  }
}

async function loadProducts(query = "") {
  if (!query) {
    query = document.getElementById("searchInput").value;
  }

  const crop = document.getElementById("searchCrop")
    ? document.getElementById("searchCrop").value
    : "";
  const region = document.getElementById("searchRegion")
    ? document.getElementById("searchRegion").value
    : "";

  const container = document.getElementById("productsContainer");
  container.innerHTML =
    '<div class="empty-state"><div class="loading"></div><p>Memuat produk...</p></div>';

  try {
    let url;
    if (query || crop || region) {
      url = `${API_BASE}/market/search?q=${encodeURIComponent(
        query
      )}&category=${encodeURIComponent(crop)}&location=${encodeURIComponent(
        region
      )}`;
    } else {
      url = `${API_BASE}/market/products`;
    }

    const res = await fetch(url);
    const data = await res.json();

    if (res.ok && data.data && data.data.length > 0) {
      container.innerHTML = data.data
        .map((p) => {
          const isOwner =
            currentUser &&
            (p.seller_id === currentUser.id || p.SellerID === currentUser.id);
          let actionButtons = "";

          if (isOwner) {
            actionButtons = `
                                <button class="btn-edit" onclick="openEditProductModal('${p.id}')"><i class="fas fa-edit"></i></button>
                                <button class="btn-delete" onclick="deleteProduct('${p.id}')"><i class="fas fa-trash"></i></button>
                            `;
          }

          return `
                        <div class="product-card">
                            <div class="product-image" onclick="viewProductDetail('${
                              p.id
                            }')">
                                ${
                                  p.image_url
                                    ? `<img src="${p.image_url}" alt="${p.name}">`
                                    : '<i class="fas fa-seedling"></i>'
                                }
                            </div>
                            <div class="product-info">
                                <h3 onclick="viewProductDetail('${p.id}')">${
            p.name
          }</h3>
                                <p>${p.description || "Tidak ada deskripsi"}</p>
                                <div class="product-footer">
                                    <div>
                                        <div class="product-price">Rp ${Number(
                                          p.price
                                        ).toLocaleString()}</div>
                                        <div class="product-stock">Stok: ${
                                          p.stock
                                        }</div>
                                    </div>
                                    <div style="display:flex; align-items:center;">
                                        <button class="btn-small" onclick="viewProductDetail('${
                                          p.id
                                        }')">Detail</button>
                                        ${actionButtons}
                                    </div>
                                </div>
                            </div>
                        </div>
                    `;
        })
        .join("");

      if (!query && !crop && !region) {
        document.getElementById("statProducts").textContent = data.data.length;
      }
    } else {
      container.innerHTML = `
                        <div class="empty-state">
                            <i class="fas fa-search" style="font-size: 3rem; margin-bottom: 1rem; opacity: 0.5;"></i>
                            <p>Produk tidak ditemukan</p>
                            ${
                              query || crop || region
                                ? `<button class="btn-secondary" onclick="resetSearch()" style="margin-top:1rem;">Tampilkan Semua</button>`
                                : ""
                            }
                        </div>`;
    }
  } catch (err) {
    container.innerHTML =
      '<div class="empty-state"><i class="fas fa-exclamation-circle"></i><p>Gagal memuat produk</p></div>';
  }
}

async function deleteProduct(id) {
  if (!confirm("Apakah Anda yakin ingin menghapus produk ini?")) return;

  try {
    const res = await fetch(`${API_BASE}/market/products/${id}`, {
      method: "DELETE",
      headers: { Authorization: `Bearer ${authToken}` },
    });

    if (res.ok) {
      showToast("Produk berhasil dihapus");
      loadProducts();
    } else {
      showToast("Gagal menghapus produk", "error");
    }
  } catch (err) {
    showToast("Terjadi kesalahan server", "error");
  }
}

async function viewProductDetail(id) {
    const modal = document.getElementById("productDetailModal");
    const content = document.getElementById("productDetailContent");
    
    content.innerHTML = '<div class="loading"></div>';
    modal.classList.add("active");

    try {
        const res = await fetch(`${API_BASE}/market/products/${id}`);
        const data = await res.json();

        if (res.ok && data.data) {
            const p = data.data;
            
            const priceFormatted = Number(p.price).toLocaleString();

            content.innerHTML = `
                <h2 style="color:var(--primary); margin-bottom:1rem;">${p.name}</h2>
                
                ${p.image_url 
                    ? `<img src="${p.image_url}" style="width:100%; max-height:300px; object-fit:cover; border-radius:10px; margin-bottom:1rem;">` 
                    : '<div style="height:200px; background:var(--dark-lighter); display:flex; align-items:center; justify-content:center; border-radius:10px; margin-bottom:1rem;"><i class="fas fa-seedling" style="font-size:3rem; color:var(--text-muted);"></i></div>'
                }
                
                <div style="display:grid; grid-template-columns:1fr 1fr; gap:1rem; margin-bottom:1rem;">
                    <div>
                        <strong>Harga:</strong> <br>
                        <span style="font-size:1.2rem; color:var(--primary); font-weight:bold;">Rp ${priceFormatted}</span>
                    </div>
                    <div>
                        <strong>Stok:</strong> <br>
                        <span>${p.stock} unit</span>
                    </div>
                    <div><strong>Kategori:</strong> <br>${p.category || "-"}</div>
                    <div><strong>Lokasi:</strong> <br>${p.location || "-"}</div>
                </div>
                
                <div style="background:var(--dark); padding:1rem; border-radius:10px; margin-bottom:1.5rem;">
                    <strong>Deskripsi:</strong>
                    <p style="margin-top:0.5rem; color:var(--text-muted); white-space:pre-line;">${p.description}</p>
                </div>

                <div style="display: flex; gap: 10px;">
                    <button class="btn btn-primary" onclick="buyProduct('${p.id}')" style="flex: 1; display:flex; align-items:center; justify-content:center; gap:8px;">
                        <i class="fas fa-shopping-cart"></i> Beli Sekarang
                    </button>
                    
                    <button class="btn" onclick="openReservationModal('${p.id}')" style="flex: 1; background: transparent; border: 2px solid var(--accent); color: var(--accent); display:flex; align-items:center; justify-content:center; gap:8px;">
                        <i class="fas fa-calendar-alt"></i> Reservasi
                    </button>
                </div>
            `;
        }
    } catch (err) {
        console.error(err);
        content.innerHTML = '<p style="text-align:center;">Gagal memuat detail produk.</p>';
    }
}

async function openEditProductModal(id) {
  try {
    const res = await fetch(`${API_BASE}/market/products/${id}`);
    const data = await res.json();

    if (res.ok && data.data) {
      const p = data.data;
      document.getElementById("editProductId").value = p.id;
      document.getElementById("editProductName").value = p.name;
      document.getElementById("editProductDescription").value = p.description;
      document.getElementById("editProductPrice").value = p.price;
      document.getElementById("editProductStock").value = p.stock;
      if (document.getElementById("editProductCrop"))
        document.getElementById("editProductCrop").value = p.category || "";
      if (document.getElementById("editProductRegion"))
        document.getElementById("editProductRegion").value = p.location || "";

      document.getElementById("editProductModal").classList.add("active");
    }
  } catch (err) {
    showToast("Gagal memuat data edit", "error");
  }
}

async function handleUpdateProduct(e) {
  e.preventDefault();
  const id = document.getElementById("editProductId").value;
  const btn = e.target.querySelector("button");
  btn.innerHTML = "Menyimpan...";
  btn.disabled = true;

  try {
    const res = await fetch(`${API_BASE}/market/products/${id}`, {
      method: "PUT",
      headers: {
        Authorization: `Bearer ${authToken}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        name: document.getElementById("editProductName").value,
        description: document.getElementById("editProductDescription").value,
        price: parseInt(document.getElementById("editProductPrice").value),
        stock: parseInt(document.getElementById("editProductStock").value),
        category: document.getElementById("editProductCrop").value,
        location: document.getElementById("editProductRegion").value,
      }),
    });

    if (res.ok) {
      showToast("Produk berhasil diperbarui");
      closeModal("editProductModal");
      loadProducts();
    } else {
      const data = await res.json();
      showToast(data.message || "Gagal update produk", "error");
    }
  } catch (err) {
    showToast("Terjadi kesalahan server", "error");
  }

  btn.innerHTML = "Simpan Perubahan";
  btn.disabled = false;
}

function handleSearch(e) {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    loadProducts();
  }, 500);
}

function resetSearch() {
  document.getElementById("searchInput").value = "";
  if (document.getElementById("searchCrop"))
    document.getElementById("searchCrop").value = "";
  if (document.getElementById("searchRegion"))
    document.getElementById("searchRegion").value = "";
  loadProducts();
}

async function handleRegister(e) {
  e.preventDefault();
  const btn = document.getElementById("registerBtn");
  btn.innerHTML = '<div class="loading"></div> Mendaftar...';
  btn.disabled = true;

  try {
    const res = await fetch(`${API_BASE}/auth/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name: document.getElementById("registerName").value,
        email: document.getElementById("registerEmail").value,
        password: document.getElementById("registerPassword").value,
      }),
    });

    const data = await res.json();

    if (res.ok) {
      showToast("Registrasi berhasil! Silakan login 🎉");
      closeModal("registerModal");
      openLoginModal();
    } else {
      showToast(data.message || "Registrasi gagal", "error");
    }
  } catch (err) {
    showToast("Gagal terhubung ke server", "error");
  }

  btn.innerHTML = "Daftar";
  btn.disabled = false;
}

async function handleLogin(e) {
  e.preventDefault();
  const btn = document.getElementById("loginBtn");
  btn.innerHTML = '<div class="loading"></div> Masuk...';
  btn.disabled = true;

  try {
    const res = await fetch(`${API_BASE}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        email: document.getElementById("loginEmail").value,
        password: document.getElementById("loginPassword").value,
      }),
    });

    const data = await res.json();

    if (res.ok && data.data.token) {
      authToken = data.data.token;
      localStorage.setItem("authToken", authToken);
      showToast("Login berhasil! Selamat datang 🌱");
      closeModal("loginModal");
      await loadUserProfile();
      loadProducts();
      loadQuestions();
    } else {
      showToast(data.message || "Login gagal", "error");
    }
  } catch (err) {
    showToast("Gagal terhubung ke server", "error");
  }

  btn.innerHTML = "Masuk";
  btn.disabled = false;
}

async function loadUserProfile() {
  try {
    const res = await fetch(`${API_BASE}/users/me`, {
      headers: { Authorization: `Bearer ${authToken}` },
    });

    if (res.ok) {
      const data = await res.json();
      currentUser = data.data;
      updateAuthUI();
    } else {
      logout();
    }
  } catch (err) {
    console.error("Failed to load profile");
  }
}

function updateAuthUI() {
  const authButtons = document.getElementById("authButtons");
  if (currentUser) {
    authButtons.innerHTML = `
            <div class="user-info">
                <button onclick="openProfile()" style="background:none; border:none; cursor:pointer; display:flex; align-items:center; gap:10px;">
                    <div style="width:35px; height:35px; background:var(--primary); border-radius:50%; display:flex; align-items:center; justify-content:center; color:white; font-weight:bold;">
                        ${currentUser.name.charAt(0).toUpperCase()}
                    </div>
                    <span class="user-name">${currentUser.name}</span>
                </button>
            </div>
        `;
    document.getElementById("createProductBtn").style.display = "block";
  } else {
    authButtons.innerHTML = `<button class="btn-login" onclick="openLoginModal()">Masuk</button>`;
    document.getElementById("createProductBtn").style.display = "none";
  }
}
async function loadMyProducts() {
  const container = document.getElementById("myProductsList");
  container.innerHTML = '<div class="loading"></div>';

  if (!currentUser) return;

  try {
    const res = await fetch(`${API_BASE}/market/products`);
    const data = await res.json();

    if (res.ok && data.data) {
      const myProds = data.data.filter(
        (p) => p.seller_id === currentUser.id || p.SellerID === currentUser.id
      );

      if (myProds.length > 0) {
        container.innerHTML = myProds
          .map(
            (p) => `
                    <div class="product-card">
                        <div class="product-image">
                             ${
                               p.image_url
                                 ? `<img src="${p.image_url}" alt="${p.name}">`
                                 : '<i class="fas fa-seedling"></i>'
                             }
                        </div>
                        <div class="product-info">
                            <h3>${p.name}</h3>
                            <p class="product-price">Rp ${Number(
                              p.price
                            ).toLocaleString()}</p>
                            <p class="product-stock">Stok: ${p.stock}</p>
                            <div style="margin-top:1rem; display:flex; gap:0.5rem;">
                                <button class="btn-edit" onclick="openEditProductModal('${
                                  p.id
                                }')">Edit</button>
                                <button class="btn-delete" onclick="deleteProduct('${
                                  p.id
                                }')">Hapus</button>
                            </div>
                        </div>
                    </div>
                `
          )
          .join("");
      } else {
        container.innerHTML =
          '<div class="empty-state"><p>Anda belum menjual produk apapun.</p></div>';
      }
    }
  } catch (err) {
    container.innerHTML = "<p>Gagal memuat produk.</p>";
  }
}

async function loadMyOrdersData() {
  const container = document.getElementById("myOrdersList");
  container.innerHTML = '<div class="loading"></div>';

  const orders = await loadMyOrders();

  if (orders && orders.length > 0) {
    container.innerHTML = orders
      .map(
        (o) => `
            <div class="order-card">
                <div>
                    <h4 style="color:var(--primary);">${
                      o.product_name || "Produk"
                    }</h4>
                    <p style="color:var(--text-muted); font-size:0.9rem;">Jumlah: ${
                      o.quantity
                    } • Total: Rp ${Number(
          o.total_price || 0
        ).toLocaleString()}</p>
                    <p style="font-size:0.8rem; margin-top:0.5rem;">${getTimeAgo(
                      o.created_at
                    )}</p>
                </div>
                <span class="badge badge-primary">Berhasil</span>
            </div>
        `
      )
      .join("");
  } else {
    container.innerHTML =
      '<div class="empty-state"><p>Belum ada riwayat pembelian.</p></div>';
  }
}

async function loadMyReservationsData() {
  const container = document.getElementById("myReservationsList");
  container.innerHTML = '<div class="loading"></div>';

  const reservations = await loadMyReservations();

  if (reservations && reservations.length > 0) {
    container.innerHTML = reservations
      .map(
        (r) => `
            <div class="order-card">
                <div>
                    <h4 style="color:var(--accent);">${
                      r.product_name || "Produk"
                    }</h4>
                    <p style="color:var(--text-muted); font-size:0.9rem;">Jumlah: ${
                      r.quantity
                    } • Catatan: ${r.note || "-"}</p>
                    <p style="font-size:0.8rem; margin-top:0.5rem;">${getTimeAgo(
                      r.created_at
                    )}</p>
                </div>
                <span class="badge badge-outline">Menunggu</span>
            </div>
        `
      )
      .join("");
  } else {
    container.innerHTML =
      '<div class="empty-state"><p>Belum ada reservasi.</p></div>';
  }
}

function logout() {
  authToken = null;
  currentUser = null;
  localStorage.removeItem("authToken");

  updateAuthUI();
  closeProfile();
  showToast("Anda telah keluar");

  loadProducts();
}

async function handleCreateProduct(e) {
  e.preventDefault();
  if (!authToken) {
    showToast("Silakan login terlebih dahulu", "error");
    return;
  }

  const btn = document.getElementById("createProductBtnSubmit");
  btn.innerHTML = '<div class="loading"></div> Menambah...';
  btn.disabled = true;

  try {
    const res = await fetch(`${API_BASE}/market/products`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${authToken}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        name: document.getElementById("productName").value,
        description: document.getElementById("productDescription").value,
        price: parseInt(document.getElementById("productPrice").value),
        stock: parseInt(document.getElementById("productStock").value),
        category: document.getElementById("productCrop").value,
        location: document.getElementById("productRegion").value,
      }),
    });

    const data = await res.json();

    if (res.ok) {
      const productId = data.data.id;
      const fileInput = document.getElementById("productImage");

      if (fileInput.files[0]) {
        await uploadProductImage(productId, fileInput.files[0]);
      }

      showToast("Produk berhasil ditambahkan! 📦");
      closeModal("createProductModal");
      loadProducts();
      e.target.reset();
    } else {
      showToast(data.message || "Gagal menambah produk", "error");
    }
  } catch (err) {
    showToast("Gagal terhubung ke server", "error");
  }

  btn.innerHTML = "Tambah Produk";
  btn.disabled = false;
}

async function uploadProductImage(productId, file) {
  const formData = new FormData();
  formData.append("image", file);

  try {
    await fetch(`${API_BASE}/market/products/${productId}/upload`, {
      method: "POST",
      headers: { Authorization: `Bearer ${authToken}` },
      body: formData,
    });
  } catch (err) {
    console.error("Failed to upload image");
  }
}

async function buyProduct(id) {
  if (!authToken) {
    showToast("Silakan login untuk membeli produk", "error");
    openLoginModal();
    return;
  }

  try {
    const res = await fetch(`${API_BASE}/market/products/${id}`);
    const data = await res.json();

    if (res.ok && data.data) {
      const p = data.data;

      if (p.stock < 1) {
        showToast("Maaf, stok produk ini habis!", "error");
        return;
      }

      document.getElementById("buyProductId").value = p.id;
      document.getElementById("buyProductName").textContent = p.name;
      document.getElementById(
        "buyProductPriceDisplay"
      ).textContent = `Rp ${Number(p.price).toLocaleString()}/unit`;
      document.getElementById("buyProductStock").textContent = p.stock;
      document.getElementById("buyProductPriceUnit").value = p.price;

      document.getElementById("buyQuantity").value = 1;
      document.getElementById("buyQuantity").max = p.stock;
      document.getElementById("buyNote").value = "";

      calculateTotal();
      closeModal("productDetailModal");
      document.getElementById("buyModal").classList.add("active");
    }
  } catch (err) {
    showToast("Gagal memuat data produk", "error");
  }
}

function calculateTotal() {
  const qtyInput = document.getElementById("buyQuantity");
  const priceUnit = parseInt(
    document.getElementById("buyProductPriceUnit").value
  );
  const stockMax = parseInt(
    document.getElementById("buyProductStock").textContent
  );

  let qty = parseInt(qtyInput.value);

  if (qty < 1) qty = 1;
  if (qty > stockMax) {
    qty = stockMax;
    showToast(`Maksimal pembelian ${stockMax} unit`, "error");
  }

  qtyInput.value = qty;

  const total = qty * priceUnit;
  document.getElementById("buyTotalPrice").textContent = `Rp ${Number(
    total
  ).toLocaleString()}`;
}

function adjustQty(change) {
  const qtyInput = document.getElementById("buyQuantity");
  let currentQty = parseInt(qtyInput.value) || 0;
  qtyInput.value = currentQty + change;
  calculateTotal();
}

async function handleOrderSubmit(e) {
  e.preventDefault();

  const btn = document.getElementById("btnConfirmOrder");
  const productId = document.getElementById("buyProductId").value;
  const quantity = parseInt(document.getElementById("buyQuantity").value);
  const note = document.getElementById("buyNote").value;

  btn.innerHTML =
    '<div class="loading" style="width:15px; height:15px;"></div> Memproses...';
  btn.disabled = true;

  try {
    const res = await fetch(`${API_BASE}/market/orders`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${authToken}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        product_id: productId,
        quantity: quantity,
        note: note,
      }),
    });

    const data = await res.json();

    if (res.ok) {
      showToast("Pesanan berhasil dibuat! 📦");
      closeModal("buyModal");

      if (typeof openProfile === "function") {
        openProfile();
        switchProfileTab("my-orders");
      } else {
        loadProducts();
      }
    } else {
      showToast(data.status || "Gagal membuat pesanan", "error");
    }
  } catch (err) {
    console.error(err);
    showToast("Terjadi kesalahan koneksi", "error");
  } finally {
    btn.innerHTML = "Buat Pesanan";
    btn.disabled = false;
  }
}

async function loadQuestions() {
  const container = document.getElementById("forumContainer");
  container.innerHTML =
    '<div class="empty-state"><div class="loading"></div><p>Memuat diskusi...</p></div>';

  try {
    const res = await fetch(`${API_BASE}/questions`, {
      headers: authToken ? { Authorization: `Bearer ${authToken}` } : {},
    });
    const data = await res.json();

    if (res.ok && data.data && data.data.length > 0) {
      container.innerHTML = data.data
        .slice(0, 5)
        .map((q) => {
          const initials = q.user?.name?.substring(0, 2).toUpperCase() || "AN";
          const timeAgo = getTimeAgo(q.created_at);

          const isLiked = q.is_liked || false;
          const heartClass = isLiked ? "fas" : "far";
          const btnClass = isLiked ? "liked" : "";

          return `
                <div class="forum-card" onclick="viewQuestion('${q.id}')">
                    <div class="forum-header">
                        <div class="forum-avatar">${initials}</div>
                        <div class="forum-meta">
                            <h4>${q.user?.name || "User"}</h4>
                            <span>${timeAgo}</span>
                        </div>
                    </div>
                    <h3>${q.title}</h3>
                    <p>${q.content.substring(0, 150)}${
            q.content.length > 150 ? "..." : ""
          }</p>
                    
                    <div class="forum-stats">
                        <div class="forum-stat">
                            <i class="fas fa-comment"></i>
                            <span>${q.answer_count || 0} Jawaban</span>
                        </div>
                        <div class="forum-stat">
                            <i class="fas fa-tag"></i>
                            <span>${q.category || "Umum"}</span>
                        </div>
                        
                        <button class="btn-like ${btnClass}" onclick="toggleLike(event, '${
            q.id
          }')" id="btn-like-${q.id}">
                            <i class="${heartClass} fa-heart"></i>
                            <span id="count-like-${q.id}">${
            q.likes_count || 0
          }</span>
                        </button>
                    </div>
                </div>
                `;
        })
        .join("");

      document.getElementById("statQuestions").textContent = data.data.length;
    } else {
      container.innerHTML =
        '<div class="empty-state"><i class="fas fa-comments"></i><p>Belum ada diskusi</p></div>';
    }
  } catch (err) {
    console.error(err);
    container.innerHTML =
      '<div class="empty-state"><i class="fas fa-exclamation-circle"></i><p>Gagal memuat diskusi</p></div>';
  }
}

async function handleAsk(e) {
  e.preventDefault();
  if (!authToken) {
    showToast("Silakan login terlebih dahulu", "error");
    closeModal("askModal");
    openLoginModal();
    return;
  }

  const btn = document.getElementById("askBtn");
  btn.innerHTML = '<div class="loading"></div> Memposting...';
  btn.disabled = true;

  try {
    const res = await fetch(`${API_BASE}/questions`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${authToken}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        title: document.getElementById("questionTitle").value,
        content: document.getElementById("questionContent").value,
        category: document.getElementById("questionCategory").value,
      }),
    });

    const data = await res.json();

    if (res.ok) {
      showToast("Pertanyaan berhasil diposting! 💬");
      closeModal("askModal");
      loadQuestions();
      e.target.reset();
    } else {
      showToast(data.message || "Gagal memposting pertanyaan", "error");
    }
  } catch (err) {
    showToast("Gagal terhubung ke server", "error");
  }

  btn.innerHTML = "Posting Pertanyaan";
  btn.disabled = false;
}

function viewQuestion(id) {
  showToast("Detail pertanyaan akan segera tersedia! 📖");
}

// ==========================================
// PEST CONTROL LOGIC
// ==========================================

async function loadAlerts() {
  const container = document.getElementById("alertsContainer");
  container.innerHTML =
    '<div class="empty-state"><div class="loading"></div><p>Memuat laporan...</p></div>';

  try {
    const res = await fetch(`${API_BASE}/alerts/map`);
    const data = await res.json();

    if (res.ok && data.data && data.data.length > 0) {
      container.innerHTML = data.data
        .slice(0, 6)
        .map((alert) => {
          const severityClass = alert.severity || "low";
          const severityText =
            {
              high: "Tinggi",
              medium: "Sedang",
              low: "Rendah",
            }[severityClass] || "Rendah";

          return `
<div class="alert-card ${severityClass}" onclick="focusMap('${alert.id}', '${
            alert.city
          }')">
                                <div class="alert-header">
                                    <h4>${alert.pest_name}</h4>
                                    <span class="severity-badge severity-${severityClass}">${severityText}</span>
                                </div>
                                <p><i class="fas fa-map-marker-alt"></i> ${
                                  alert.city
                                }</p>
                                <p>${alert.description}</p>
                                <div style="display: flex; justify-content: space-between; align-items: center; margin-top: 1rem;">
                                    <span style="color: var(--text-muted); font-size: 0.85rem;">
                                        <i class="fas fa-check-circle"></i> ${
                                          alert.verification_count || 0
                                        } Verifikasi
                                    </span>
                                    <span style="color: var(--text-muted); font-size: 0.85rem;">
                                        ${getTimeAgo(alert.created_at)}
                                    </span>
                                </div>
                            </div>
                        `;
        })
        .join("");
      document.getElementById("statAlerts").textContent = data.data.length;
    } else {
      container.innerHTML =
        '<div class="empty-state"><i class="fas fa-bug"></i><p>Belum ada laporan hama</p></div>';
    }
  } catch (err) {
    container.innerHTML =
      '<div class="empty-state"><i class="fas fa-exclamation-circle"></i><p>Gagal memuat laporan</p></div>';
  }
}

async function handleReport(e) {
  e.preventDefault();
  if (!authToken) {
    showToast("Silakan login terlebih dahulu", "error");
    closeModal("reportModal");
    openLoginModal();
    return;
  }

  const btn = document.getElementById("reportBtn");
  btn.innerHTML = '<div class="loading"></div> Mengirim...';
  btn.disabled = true;

  try {
    const res = await fetch(`${API_BASE}/alerts`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${authToken}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        pest_name: document.getElementById("pestName").value,
        description: document.getElementById("pestDescription").value,
        city: document.getElementById("pestCity").value,
        severity: document.getElementById("pestSeverity").value,
      }),
    });

    const data = await res.json();

    if (res.ok) {
      showToast("Laporan hama berhasil dikirim! 🐛");
      closeModal("reportModal");
      loadAlerts();
      e.target.reset();
    } else {
      showToast(data.message || "Gagal mengirim laporan", "error");
    }
  } catch (err) {
    showToast("Gagal terhubung ke server", "error");
  }

  btn.innerHTML = "Kirim Laporan";
  btn.disabled = false;
}

async function viewAlertDetail(id) {
  if (!id) return;

  try {
    const res = await fetch(`${API_BASE}/alerts/${id}`);
    const data = await res.json();

    if (res.ok && data.data) {
      const alert = data.data;
      const severityText =
        {
          high: "Tinggi",
          medium: "Sedang",
          low: "Rendah",
        }[alert.severity] || "Rendah";

      showToast(
        `${alert.pest_name} di ${alert.city} - Tingkat: ${severityText}`
      );
    }
  } catch (err) {
    console.error("Error loading alert detail:", err);
  }
}

// ==========================================
// UTILITY & HELPER FUNCTIONS
// ==========================================

async function createReservation(productId, quantity, note) {
  if (!authToken) {
    showToast("Silakan login terlebih dahulu", "error");
    openLoginModal();
    return null;
  }

  try {
    const res = await fetch(`${API_BASE}/market/reservations`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${authToken}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        product_id: productId,
        quantity: quantity,
        note: note || "",
      }),
    });

    const data = await res.json();

    if (res.ok) {
      showToast("Reservasi berhasil dibuat! 📦");
      return data.data;
    } else {
      showToast(data.message || "Gagal membuat reservasi", "error");
      return null;
    }
  } catch (err) {
    showToast("Gagal terhubung ke server", "error");
    return null;
  }
}

async function createOrder(productId, quantity, note) {
  if (!authToken) {
    showToast("Silakan login terlebih dahulu", "error");
    openLoginModal();
    return null;
  }

  try {
    const res = await fetch(`${API_BASE}/market/orders`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${authToken}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        product_id: productId,
        quantity: quantity,
        note: note || "",
      }),
    });

    const data = await res.json();

    if (res.ok) {
      showToast("Pesanan berhasil dibuat! 🛒");
      return data.data;
    } else {
      showToast(data.message || "Gagal membuat pesanan", "error");
      return null;
    }
  } catch (err) {
    showToast("Gagal terhubung ke server", "error");
    return null;
  }
}

async function loadMyReservations() {
  if (!authToken) return;

  try {
    const res = await fetch(`${API_BASE}/users/me/reservations`, {
      headers: { Authorization: `Bearer ${authToken}` },
    });

    const data = await res.json();

    if (res.ok && data.data) {
      console.log("My Reservations:", data.data);
      return data.data;
    }
  } catch (err) {
    console.error("Error loading reservations:", err);
  }
  return [];
}

async function loadMyOrders() {
  if (!authToken) return;

  try {
    const res = await fetch(`${API_BASE}/users/me/orders`, {
      headers: { Authorization: `Bearer ${authToken}` },
    });

    const data = await res.json();

    if (res.ok && data.data) {
      console.log("My Orders:", data.data);
      return data.data;
    }
  } catch (err) {
    console.error("Error loading orders:", err);
  }
  return [];
}

async function submitAnswer(questionId, content) {
  if (!authToken) {
    showToast("Silakan login terlebih dahulu", "error");
    openLoginModal();
    return null;
  }

  try {
    const res = await fetch(`${API_BASE}/questions/${questionId}/answers`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${authToken}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ content }),
    });

    const data = await res.json();

    if (res.ok) {
      showToast("Jawaban berhasil dikirim! 💬");
      return data.data;
    } else {
      showToast(data.message || "Gagal mengirim jawaban", "error");
      return null;
    }
  } catch (err) {
    showToast("Gagal terhubung ke server", "error");
    return null;
  }
}

async function toggleLike(event, questionId) {
  event.stopPropagation();

  if (!authToken) {
    showToast("Silakan login dulu bosku! 🤭", "error");
    openLoginModal();
    return;
  }

  const btn = document.getElementById(`btn-like-${questionId}`);

  const targetBtn = btn || event.currentTarget;
  const icon = targetBtn.querySelector("i");
  const countSpan = targetBtn.querySelector("span");

  let currentCount = parseInt(countSpan.innerText) || 0;

  try {
    const res = await fetch(`${API_BASE}/questions/${questionId}/like`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${authToken}`,
        "Content-Type": "application/json",
      },
    });

    if (res.ok) {
      const isLikedNow = targetBtn.classList.contains("liked");

      if (isLikedNow) {
        targetBtn.classList.remove("liked");
        icon.classList.remove("fas");
        icon.classList.add("far");
        countSpan.innerText = Math.max(0, currentCount - 1);
      } else {
        targetBtn.classList.add("liked");
        icon.classList.remove("far");
        icon.classList.add("fas");
        countSpan.innerText = currentCount + 1;

        icon.style.transform = "scale(1.3)";
        setTimeout(() => (icon.style.transform = "scale(1)"), 200);
      }
    } else {
      const data = await res.json();
      showToast(data.message || "Gagal like", "error");
    }
  } catch (err) {
    console.error(err);
    showToast("Gagal koneksi server", "error");
  }
}

function getTimeAgo(dateString) {
  if (!dateString) return "Baru saja";

  const date = new Date(dateString);
  const now = new Date();
  const seconds = Math.floor((now - date) / 1000);

  const intervals = {
    tahun: 31536000,
    bulan: 2592000,
    minggu: 604800,
    hari: 86400,
    jam: 3600,
    menit: 60,
  };

  for (const [unit, secondsInUnit] of Object.entries(intervals)) {
    const interval = Math.floor(seconds / secondsInUnit);
    if (interval >= 1) {
      return `${interval} ${unit} yang lalu`;
    }
  }

  return "Baru saja";
}

function showToast(message, type = "success") {
  const toast = document.getElementById("toast");
  toast.textContent = message;
  toast.className = `toast show ${type}`;

  setTimeout(() => {
    toast.classList.remove("show");
  }, 3000);
}

function openLoginModal() {
  document.getElementById("loginModal").classList.add("active");
}

function openRegisterModal() {
  document.getElementById("registerModal").classList.add("active");
}

function openCreateProductModal() {
  if (!authToken) {
    showToast("Silakan login terlebih dahulu", "error");
    openLoginModal();
    return;
  }
  document.getElementById("createProductModal").classList.add("active");
}

function openAskModal() {
  if (!authToken) {
    showToast("Silakan login terlebih dahulu", "error");
    openLoginModal();
    return;
  }
  document.getElementById("askModal").classList.add("active");
}

async function viewQuestion(id) {
  const modal = document.getElementById("questionDetailModal");
  const content = document.getElementById("questionDetailContent");

  // Tampilkan loading state
  content.innerHTML =
    '<div class="empty-state"><div class="loading"></div><p>Memuat diskusi...</p></div>';
  modal.classList.add("active");

  try {
    // Fetch data pertanyaan detail
    const res = await fetch(`${API_BASE}/questions/${id}`);
    const data = await res.json();

    if (res.ok && data.data) {
      const q = data.data;
      const initials = q.user?.name?.substring(0, 2).toUpperCase() || "AN";
      const timeAgo = getTimeAgo(q.created_at);

      // Render Jawaban (Jika ada array answers, jika tidak kosongkan)
      // Asumsi backend mengembalikan array 'answers' di dalam objek detail
      const answersHtml =
        (q.answers || [])
          .map(
            (ans) => `
                <div class="answer-card">
                    <div class="answer-header">
                        <span class="answer-author">${
                          ans.user?.name || "User"
                        }</span>
                        <span class="answer-date">${getTimeAgo(
                          ans.created_at
                        )}</span>
                    </div>
                    <div class="answer-content">${ans.content}</div>
                </div>
            `
          )
          .join("") ||
        '<p style="text-align:center; color:var(--text-muted);">Belum ada jawaban. Jadilah yang pertama menjawab!</p>';

      // Render HTML Lengkap
      content.innerHTML = `
                <div class="question-full-header">
                    <div style="display:flex; align-items:center; gap:1rem; margin-bottom:1rem;">
                        <div class="forum-avatar" style="width:40px; height:40px; font-size:1rem;">${initials}</div>
                        <div>
                            <div style="font-weight:600; color:var(--text);">${
                              q.user?.name || "Anonymous"
                            }</div>
                            <div style="font-size:0.85rem; color:var(--text-muted);">${timeAgo}</div>
                        </div>
                    </div>
                    <h2 class="question-full-title">${q.title}</h2>
                    <div>
                        <span class="question-meta-badge"><i class="fas fa-tag"></i> ${
                          q.category || "Umum"
                        }</span>
                        <span class="question-meta-badge" style="cursor:pointer;" onclick="toggleLike('${
                          q.id
                        }')">
                            <i class="fas fa-heart"></i> ${
                              q.likes_count || 0
                            } Likes
                        </span>
                    </div>
                </div>

                <div class="question-body">${q.content}</div>

                <div class="answers-section">
                    <div class="answers-header">
                        <i class="fas fa-comments"></i> Jawaban (${
                          q.answers ? q.answers.length : 0
                        })
                    </div>
                    <div id="answersList">
                        ${answersHtml}
                    </div>
                </div>

                <div class="reply-box">
                    <h4 style="margin-bottom:1rem; color:var(--text);">Berikan Jawaban Anda</h4>
                    <form onsubmit="handleReplySubmit(event, '${q.id}')">
                        <textarea class="reply-textarea" id="replyContent" placeholder="Tulis solusi atau saran Anda di sini..." required></textarea>
                        <button type="submit" class="btn btn-primary" style="float:right;">
                            <i class="fas fa-paper-plane"></i> Kirim Jawaban
                        </button>
                        <div style="clear:both;"></div>
                    </form>
                </div>
            `;
    } else {
      content.innerHTML =
        '<div class="empty-state"><p>Gagal memuat data pertanyaan.</p></div>';
    }
  } catch (err) {
    console.error(err);
    content.innerHTML =
      '<div class="empty-state"><p>Terjadi kesalahan koneksi.</p></div>';
  }
}

function openProfile() {
  document.getElementById("home").style.display = "none";
  document.getElementById("features").style.display = "none";
  document.getElementById("products").style.display = "none";
  document.getElementById("forum").style.display = "none";
  document.getElementById("pest").style.display = "none";
  document.querySelector(".stats").style.display = "none";
  document.querySelector(".hero").style.display = "none";
  document.getElementById("profileSection").style.display = "block";

  if (currentUser) {
    document.getElementById("profileName").textContent = currentUser.name;
    document.getElementById("profileEmail").textContent = currentUser.email;
  }
  loadMyProducts();
  window.scrollTo(0, 0);
}

function switchProfileTab(tabName) {
  document
    .querySelectorAll(".tab-btn")
    .forEach((btn) => btn.classList.remove("active"));
  document
    .querySelectorAll(".profile-tab-content")
    .forEach((content) => (content.style.display = "none"));

  const contentTarget = document.getElementById(`tab-${tabName}`);
  if (contentTarget) {
    contentTarget.style.display = "block";
  }

  const btnTarget = document.getElementById(`tab-btn-${tabName}`);
  if (btnTarget) {
    btnTarget.classList.add("active");
  }

  if (tabName === "my-products") loadMyProducts();
  if (tabName === "my-orders") loadMyOrdersData();
  if (tabName === "my-reservations") loadMyReservationsData();
}

function closeProfile() {
  document.getElementById("home").style.display = "flex";
  document.getElementById("features").style.display = "block";
  document.getElementById("products").style.display = "block";
  document.getElementById("forum").style.display = "block";
  document.getElementById("pest").style.display = "block";
  document.querySelector(".stats").style.display = "block";
  document.querySelector(".hero").style.display = "flex";

  document.getElementById("profileSection").style.display = "none";
  window.scrollTo(0, 0);
}

function openEditProfileModal() {
  if (!currentUser) return;

  document.getElementById("editProfileName").value = currentUser.name;
  document.getElementById("editProfileEmail").value = currentUser.email;
  document.getElementById("editProfilePassword").value = "";
  document.getElementById("editProfileModal").classList.add("active");
}

async function handleUpdateProfile(e) {
  e.preventDefault();

  const btn = document.getElementById("btnUpdateProfile");
  const nameInput = document.getElementById("editProfileName").value;
  const passwordInput = document.getElementById("editProfilePassword").value;

  if (!nameInput) {
    showToast("Nama tidak boleh kosong", "error");
    return;
  }

  const payload = {
    name: nameInput,
  };

  if (passwordInput) {
    if (passwordInput.length < 8) {
      showToast("Password minimal 8 karakter", "error");
      return;
    }
    payload.password = passwordInput;
  }

  btn.innerHTML =
    '<div class="loading" style="width:15px; height:15px;"></div> Menyimpan...';
  btn.disabled = true;

  try {
    const res = await fetch(`${API_BASE}/users/me`, {
      method: "PATCH",
      headers: {
        Authorization: `Bearer ${authToken}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });

    const data = await res.json();

    if (res.ok) {
      currentUser.name = nameInput;
      document.getElementById("profileName").textContent = nameInput;
      updateAuthUI();

      showToast("Profil berhasil diperbarui! 🎉");
      closeModal("editProfileModal");
    } else {
      showToast(data.message || "Gagal update profil", "error");
    }
  } catch (err) {
    console.error(err);
    showToast("Terjadi kesalahan koneksi", "error");
  } finally {
    btn.innerHTML = "Simpan Perubahan";
    btn.disabled = false;
  }
}

async function handleReplySubmit(e, questionId) {
  e.preventDefault();

  if (!authToken) {
    showToast("Silakan login untuk menjawab", "error");
    openLoginModal();
    return;
  }

  const textarea = document.getElementById("replyContent");
  const content = textarea.value;
  const btn = e.target.querySelector("button");

  const originalBtnText = btn.innerHTML;
  btn.innerHTML =
    '<div class="loading" style="width:15px; height:15px; border-width:2px;"></div> Mengirim...';
  btn.disabled = true;

  try {
    const result = await submitAnswer(questionId, content);

    if (result) {
      textarea.value = "";
      viewQuestion(questionId);
    }
  } catch (err) {
    showToast("Gagal mengirim jawaban", "error");
  } finally {
    btn.innerHTML = originalBtnText;
    btn.disabled = false;
  }
}

async function handleVerify(event, alertId) {
  event.stopPropagation();

  if (!authToken) {
    showToast("Login dulu untuk memverifikasi laporan", "error");
    openLoginModal();
    return;
  }

  const btn = event.currentTarget;
  const countSpan = btn.querySelector("span");

  const isCurrentlyVerified = btn.classList.contains("verified");
  let currentCount = parseInt(countSpan.innerText) || 0;

  btn.style.opacity = "0.7";

  try {
    const res = await fetch(`${API_BASE}/alerts/${alertId}/verify`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${authToken}`,
        "Content-Type": "application/json",
      },
    });

    const data = await res.json();

    if (res.ok) {
      if (isCurrentlyVerified) {
        btn.classList.remove("verified");
        countSpan.innerText = Math.max(0, currentCount - 1);
        showToast("Batal memverifikasi");
      } else {
        btn.classList.add("verified");
        countSpan.innerText = currentCount + 1;
        showToast("Laporan terverifikasi valid ✅");

        btn.style.transform = "scale(1.1)";
        setTimeout(() => (btn.style.transform = "scale(1)"), 200);
      }
    } else {
      showToast(data.message || "Gagal verifikasi", "error");
    }
  } catch (err) {
    console.error(err);
    showToast("Gagal koneksi server", "error");
  } finally {
    btn.style.opacity = "1";
  }
}

function openReportModal() {
  if (!authToken) {
    showToast("Silakan login terlebih dahulu", "error");
    openLoginModal();
    return;
  }
  document.getElementById("reportModal").classList.add("active");
}

function closeModal(modalId) {
  document.getElementById(modalId).classList.remove("active");
}

document.addEventListener("click", (e) => {
  if (e.target.classList.contains("modal")) {
    e.target.classList.remove("active");
  }
});

function animateCounter(elementId, target) {
  const element = document.getElementById(elementId);
  if (!element) return;

  let current = 0;
  const increment = target / 50;
  const timer = setInterval(() => {
    current += increment;
    if (current >= target) {
      element.textContent = target;
      clearInterval(timer);
    } else {
      element.textContent = Math.floor(current);
    }
  }, 30);
}

const observerOptions = {
  threshold: 0.5,
};

const observer = new IntersectionObserver((entries) => {
  entries.forEach((entry) => {
    if (entry.isIntersecting) {
      const statsSection = entry.target;
      animateCounter("statUsers", Math.floor(Math.random() * 5000) + 1000);
      observer.unobserve(statsSection);
    }
  });
}, observerOptions);

const statsSection = document.querySelector(".stats");
if (statsSection) {
  observer.observe(statsSection);
}

window.AgroHub = {
  createReservation:
    typeof createReservation !== "undefined" ? createReservation : undefined,
  createOrder: typeof createOrder !== "undefined" ? createOrder : undefined,
  loadMyReservations:
    typeof loadMyReservations !== "undefined" ? loadMyReservations : undefined,
  loadMyOrders: typeof loadMyOrders !== "undefined" ? loadMyOrders : undefined,
  submitAnswer: typeof submitAnswer !== "undefined" ? submitAnswer : undefined,
  toggleLike: typeof toggleLike !== "undefined" ? toggleLike : undefined,
  viewAlertDetail:
    typeof viewAlertDetail !== "undefined" ? viewAlertDetail : undefined,
};

console.log("✅ AgroHub System Loaded Successfully");

let map;
let markers = {};

const cityCoordinates = {
  Bangkalan: [-7.03, 112.924],
  Banyuwangi: [-8.219, 114.369],
  Blitar: [-8.095, 112.16],
  Bojonegoro: [-7.15, 111.881],
  Bondowoso: [-7.913, 113.821],
  Gresik: [-7.155, 112.656],
  Jember: [-8.172, 113.7],
  Jombang: [-7.545, 112.233],
  Kediri: [-7.848, 112.017],
  Lamongan: [-7.128, 112.313],
  Lumajang: [-8.133, 113.222],
  Madiun: [-7.629, 111.523],
  Magetan: [-7.653, 111.328],
  Malang: [-7.979, 112.63],
  Mojokerto: [-7.472, 112.433],
  Nganjuk: [-7.603, 111.901],
  Ngawi: [-7.403, 111.444],
  Pacitan: [-8.196, 111.106],
  Pamekasan: [-7.161, 113.483],
  Pasuruan: [-7.645, 112.907],
  Ponorogo: [-7.872, 111.462],
  Probolinggo: [-7.754, 113.215],
  Sampang: [-7.187, 113.243],
  Sidoarjo: [-7.447, 112.718],
  Situbondo: [-7.706, 114.004],
  Sumenep: [-7.008, 113.859],
  Trenggalek: [-8.05, 111.716],
  Tuban: [-6.897, 112.064],
  Tulungagung: [-8.077, 111.9],
  "Kota Batu": [-7.867, 112.526],
  "Kota Surabaya": [-7.257, 112.752],
};

function initMap() {
  if (!document.getElementById("pestMap")) return;

  map = L.map("pestMap").setView([-7.7, 112.5], 8);

  L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
    attribution: "&copy; OpenStreetMap contributors",
  }).addTo(map);
}

async function loadAlerts() {
  const container = document.getElementById("alertsContainer");
  container.innerHTML =
    '<div class="empty-state"><div class="loading"></div><p>Memuat laporan...</p></div>';
  if (!map) initMap();

  try {
    const res = await fetch(`${API_BASE}/alerts/map`);
    const data = await res.json();

    if (res.ok && data.data && data.data.length > 0) {
      container.innerHTML = data.data
        .map((alert) => {
          const severityClass = alert.severity || "low";
          const severityText =
            { high: "Tinggi", medium: "Sedang", low: "Rendah" }[
              severityClass
            ] || "Rendah";

          addMarkerToMap(alert);

          return `
    <div class="alert-card ${severityClass}" onclick="focusMap('${
            alert.id
          }', '${alert.city}')" id="card-${alert.id}">
        <div class="alert-header">
            <h4>${alert.pest_name}</h4>
            <span class="severity-badge severity-${severityClass}">${severityText}</span>
        </div>
        <p><i class="fas fa-map-marker-alt"></i> ${alert.city}</p>
        <p style="font-size:0.9rem; color:var(--text-muted); margin-bottom: 1rem;">
            ${alert.description.substring(0, 80)}...
        </p>
        
        <div style="display: flex; justify-content: space-between; align-items: center; border-top: 1px solid rgba(255,255,255,0.05); padding-top: 10px;">
            <span style="color: var(--text-muted); font-size: 0.85rem;">
                ${getTimeAgo(alert.created_at)}
            </span>

            <button 
                class="btn-verify ${alert.is_verified ? "verified" : ""}" 
                onclick="handleVerify(event, '${alert.id}')"
                id="btn-verif-${alert.id}"
            >
                <i class="fas fa-check-circle"></i> 
                <span id="count-verif-${alert.id}">${
            alert.verification_count || 0
          }</span> Valid
            </button>
        </div>
    </div>
`;
        })
        .join("");

      document.getElementById("statAlerts").textContent = data.data.length;
    } else {
      container.innerHTML =
        '<div class="empty-state"><i class="fas fa-bug"></i><p>Belum ada laporan hama</p></div>';
    }
  } catch (err) {
    console.error(err);
    container.innerHTML =
      '<div class="empty-state"><i class="fas fa-exclamation-circle"></i><p>Gagal memuat laporan</p></div>';
  }
}

function addMarkerToMap(alert) {
  const coords = cityCoordinates[alert.city];

  if (coords) {
    const colorMap = {
      high: "#ef4444",
      medium: "#f59e0b",
      low: "#10b981",
    };

    const marker = L.circle(coords, {
      color: colorMap[alert.severity] || "blue",
      fillColor: colorMap[alert.severity] || "blue",
      fillOpacity: 0.5,
      radius: 4000,
    }).addTo(map);

    const popupContent = `
            <div style="min-width: 200px;">
                <h3 style="margin:0; color:${colorMap[alert.severity]};">${
      alert.pest_name
    }</h3>
                <small style="color:#666;"><i class="fas fa-map-marker-alt"></i> ${
                  alert.city
                }</small>
                
                <p style="margin: 10px 0; font-size:0.9rem; line-height:1.4;">
                    "${alert.description}"
                </p>
                
                <div style="border-top:1px solid #eee; padding-top:8px; display:flex; justify-content:space-between; align-items:center;">
                    <span style="font-weight:bold; font-size:0.8rem;">Status: ${alert.severity.toUpperCase()}</span>
                    <span style="color:var(--primary); font-size:0.8rem;">
                        <i class="fas fa-check-circle"></i> ${
                          alert.verification_count || 0
                        } Verifikasi
                    </span>
                </div>
            </div>
        `;

    marker.bindPopup(popupContent);
    markers[alert.id] = marker;
  }
}

function focusMap(id, cityName) {
  const coords = cityCoordinates[cityName];
  if (coords && map) {
    map.flyTo(coords, 12, {
      duration: 1.5,
    });
    document
      .getElementById("pestMap")
      .scrollIntoView({ behavior: "smooth", block: "center" });
  }
  const marker = markers[id];
  if (marker) {
    setTimeout(() => {
      marker.openPopup();
    }, 800);
  }
}
