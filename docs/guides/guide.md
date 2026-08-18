# QUY TRÌNH TỰ TÌM HIỂU MỘT HỆ THỐNG TỪ ĐẦU 

Dành cho việc chinh phục những hệ thống chưa biết 

Quy trình 8 bước 
1. Product Understanding 
### Người dùng có thể làm gì 

### Hệ thống phải làm gì 
2. Black-box Explaination
- Sử dụng sản phẩm như một người dùng 
- Hiểu hệ thống từ bên ngoài trước khi nhìn vào bên trong 
- Tiến hành nhập input từng chức năng và ghi ra output   
3. Reverse Engineering
- Nếu như có source code của dự án. Có thể tiến hành tìm hiểu về request flow: Request của người dùng sẽ đi qua những gì  
4. Architecture Modeling 
- Tự vẽ Architecture 
- So sánh với các mô hình mẫu và hỏi TẠI SAO 
5. Component Design 
- Với mỗi Component bên trong Design Architecture, hãy hỏi 5 câu: 
    + Nó làm gì? 
    + Nó giao tiếp với ai? 
    + Input / Output là gì? 
    + nếu bỏ nó thì chuyện gì xảy ra? 
6. Theory on Demand 
- Trong quá trình xử lý kiến trúc, thấy cần kiến thức gì thì học kiến thức đó. 
7. Implementating 
- Tiến hành implement theo từng version, với từng cải tiến: Hello Agent -> Tool using -> Agent loop -> Session -> Multiple Tools -> Gateway -> Multi agent 
- Hãy hỏi câu hỏi: Minimal system cần những gì để tạo ra 1 hoàn chỉnh nhỏ nhất? 
8. Testing + Refelectation
9. Quay lại 3,4,5 
