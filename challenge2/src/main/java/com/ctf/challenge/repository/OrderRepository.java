package com.ctf.challenge.repository;

import com.ctf.challenge.model.Order;
import com.ctf.challenge.model.OrderItem;
import com.ctf.challenge.model.Product;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.core.RowMapper;
import org.springframework.jdbc.support.GeneratedKeyHolder;
import org.springframework.jdbc.support.KeyHolder;
import org.springframework.stereotype.Repository;

import java.sql.PreparedStatement;
import java.sql.Statement;
import java.util.List;

@Repository
public class OrderRepository {

    @Autowired
    private JdbcTemplate jdbcTemplate;

    @Autowired
    private ProductRepository productRepository;

    private final RowMapper<Order> orderMapper = (rs, rowNum) -> {
        Order o = new Order();
        o.setId(rs.getInt("id"));
        o.setUserId(rs.getInt("user_id"));
        o.setStatus(rs.getString("status"));
        o.setTotalPrice(rs.getBigDecimal("total_price"));
        o.setCreatedAt(rs.getTimestamp("created_at"));
        return o;
    };

    private final RowMapper<OrderItem> orderItemMapper = (rs, rowNum) -> {
        OrderItem item = new OrderItem();
        item.setId(rs.getInt("id"));
        item.setOrderId(rs.getInt("order_id"));
        item.setProductId(rs.getInt("product_id"));
        item.setQuantity(rs.getInt("quantity"));
        item.setPrice(rs.getBigDecimal("price"));

        Product p = productRepository.findById(item.getProductId());
        item.setProduct(p);
        return item;
    };

    public Integer save(Order order) {
        KeyHolder keyHolder = new GeneratedKeyHolder();
        jdbcTemplate.update(connection -> {
            PreparedStatement ps = connection.prepareStatement(
                    "INSERT INTO orders (user_id, status, total_price) VALUES (?, ?, ?)",
                    Statement.RETURN_GENERATED_KEYS);
            ps.setInt(1, order.getUserId());
            ps.setString(2, order.getStatus());
            ps.setBigDecimal(3, order.getTotalPrice());
            return ps;
        }, keyHolder);
        return keyHolder.getKey().intValue();
    }

    public void saveItem(OrderItem item) {
        jdbcTemplate.update("INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)",
                item.getOrderId(), item.getProductId(), item.getQuantity(), item.getPrice());
    }

    public List<Order> findByUserId(Integer userId) {
        return jdbcTemplate.query("SELECT * FROM orders WHERE user_id = ? ORDER BY created_at DESC", orderMapper,
                userId);
    }

    public List<Order> findAll() {
        return jdbcTemplate.query("SELECT * FROM orders ORDER BY created_at DESC", orderMapper);
    }

    public Order findById(Integer orderId) {
        try {
            return jdbcTemplate.queryForObject("SELECT * FROM orders WHERE id = ?", orderMapper, orderId);
        } catch (Exception e) {
            return null;
        }
    }

    public List<OrderItem> findItemsByOrderId(Integer orderId) {
        return jdbcTemplate.query("SELECT * FROM order_items WHERE order_id = ?", orderItemMapper, orderId);
    }

    public void updateStatus(Integer orderId, String status) {
        jdbcTemplate.update("UPDATE orders SET status = ? WHERE id = ?", status, orderId);
    }
}
