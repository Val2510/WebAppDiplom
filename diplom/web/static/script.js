let userEmail = '';
let userCity = '';
let userName = '';

async function registerUser() {
    const registrationForm = $('#registrationForm');
    const formData = registrationForm.serializeArray();
    const data = {};
    $(formData).each(function(_index, obj){
        data[obj.name] = obj.value;
    });

    try {
        $.ajax({
            url: '/register_user',
            type: 'POST',
            contentType: 'application/json',
            data: JSON.stringify(data),
            success: function () {
                alert('Registration successful!');
                localStorage.setItem('userEmail', data.email);
                $('#registrationModal').hide();
                registrationForm[0].reset();
                userEmail = data.email;
                userName = data.name;
                userCity = data.city;
            },
            error: function (_xhr, _status, error) {
                alert('Registration failed: ' + error.responseText);
            }
        });
    } catch (error) {
        console.error('Error during registration:', error.responseText);
        alert('Registration failed: An error occurred.');
    }
}

async function fetchUserInfo() {
    try {
        $.ajax({
            url: `/api/userinfo?email=${userEmail}`,
            method: "GET",
            dataType: "json",
            success: function (user) {
                userEmail = user.email; 
                userName = user.name;
                userCity = user.city;
                getBooks(); 
            },
            error: function (_xhr, _status, error) {
                console.error('Error fetching user info:', error.responseText);
            }
        });
    } catch (error) {
        console.error('Error fetching user info:', error.responseText);
    }
}



async function addBook() {
    const form = document.getElementById('addBookForm');
    const formData = new FormData(form);
    const bookData = {
        author: formData.get('author'),
        title: formData.get('title'),
        email: userEmail,
        city: formData.get('city'),
        name: formData.get('name'),
    };

    try {
        $.ajax({
            url: "/api/add_book",
            method: "POST",
            dataType: "json",
            contentType: "application/json",
            data: JSON.stringify(bookData),
            success: function (response) {
                form.reset();
                document.getElementById('addBookModal').style.display = 'none';
                $("#book-list1-body").empty();
                response.forEach(function(book) {
                    var row = "<tr>" +
                                "<td>" + (book.title || "") + "</td>" +
                                "<td>" + (book.author || "") + "</td>" +
                                "<td>" + (book.name ? book.name + " (" + book.email + ")" : "Unknown") + "</td>" +
                              "</tr>";
                    $("#book-list1-body").append(row);
                });
            },
            error: function (_xhr, _status, error) {
                console.error(`Failed to add book: ${error}`);
            }
        });
    } catch (error) {
        console.error('Error adding book:', error.responseText);
    }
}

function getBooks() {
    $("#book-list-body").show();
    try {
        $.ajax({
            url: "/books",
            method: "GET",
            dataType: "json", 
            success: function(response) {
                $("#book-list-body").empty();
                response.forEach(function(book) {
                    var row = "<tr>" +
                                "<td>" + book.Title + "</td>" +
                                "<td>" + book.Author + "</td>" +
                                "<td>" + (book.name ? book.name + " (" + book.email + ")" : "Unknown") + "</td>" +
                              "</tr>";
                    $("#book-list-body").append(row);
                });
            },
            error: function(_xhr, _status, error) {
                console.error("Failed to get books:", error.responseText);
            }
        });
    } catch (error) {
        console.error('Error fetching user info:', error.responseText);
    }
}

function searchBooks() {
    var query = $("#searchInput").val();
    console.log("Search query:", query); 
    $.ajax({
        url: "/search?query=" + query,
        method: "GET",
        dataType: "json",
        success: function(response) {
            console.log("Received response:", response); 

            if (response && Array.isArray(response)) {
                console.log("Processing books:", response); 
                
                $("#book-list-body").empty();
                response.forEach(function(book) {
                    var row = "<tr>" +
                                "<td>" + (book.title || "") + "</td>" +
                                "<td>" + (book.author || "") + "</td>" +
                                "<td>" + (book.name ? book.name + " (" + book.email + ")" : "Unknown") + "</td>" +
                              "</tr>";
                    $("#book-list-body").append(row);
                });
            } else {
                console.error("Response is null or not an array:", response);
            }
        },
        error: function(_xhr, _status, error) {
            console.error("Failed to search books:", error.responseText);
        }
    });
}

async function showOwnerInfo(email) {
    try {
        $.ajax({
            url: `/owner?email=${email}`,
            method: "GET",
            dataType: "json",
            success: function (owner) {
                const ownerInfo = document.getElementById('ownerInfo');
                ownerInfo.innerHTML = `Name: ${owner.Name}<br>Email: ${owner.Email}`;
                document.getElementById('ownerInfoModal').style.display = 'block';
            },
            error: function (_xhr, _status, error) {
                console.error('Error fetching owner info:', error.responseText);
            }
        });
    } catch (error) {
        console.error('Error fetching owner info:', error.responseText);
    }
}

function closeModal() {
    document.getElementById('ownerInfoModal').style.display = 'none';
}

document.addEventListener('DOMContentLoaded', async () => {
    userEmail = localStorage.getItem('userEmail');
    if (!userEmail) {
        alert('User email not found. Please register or log in.');
        return;
    }

    const registrationButton = document.getElementById('registrationButton');
    const registrationModal = document.getElementById('registrationModal');
    const closeSpan = document.querySelector('.modal .close');

    if (registrationButton && registrationModal && closeSpan) {
        registrationButton.addEventListener('click', () => {
            registrationModal.style.display = 'block';
        });

        closeSpan.addEventListener('click', () => {
            registrationModal.style.display = 'none';
        });

        window.addEventListener('click', (event) => {
            if (event.target == registrationModal) {
                registrationModal.style.display = 'none';
            }
        });

        const registrationForm = document.getElementById('registrationForm');
        registrationForm.addEventListener('submit', async (event) => {
            event.preventDefault();
            registerUser();
        });
    }

    const addBookButton = document.getElementById('add-book-button');
    const closeAddBookModal = document.getElementById('close-add-book-modal');
    const addBookForm = document.getElementById('addBookForm');

    if (addBookButton && closeAddBookModal && addBookForm) {
        addBookButton.addEventListener('click', () => {
            document.getElementById('addBookModal').style.display = 'block';
        });

        closeAddBookModal.addEventListener('click', () => {
            document.getElementById('addBookModal').style.display = 'none';
        });

        window.addEventListener('click', (event) => {
            if (event.target === document.getElementById('addBookModal')) {
                document.getElementById('addBookModal').style.display = 'none';
            }
        });

        addBookForm.addEventListener('submit', async (event) => {
            event.preventDefault();
            addBook();
        });
    }
});
