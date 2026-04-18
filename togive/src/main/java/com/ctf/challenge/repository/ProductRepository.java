package com.ctf.challenge.repository;

import com.ctf.challenge.model.Product;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.core.RowMapper;
import org.springframework.stereotype.Repository;
import java.util.List;

@Repository
public class ProductRepository {

    @Autowired
    private JdbcTemplate jdbcTemplate;

    private final RowMapper<Product> rowMapper = (rs, rowNum) -> {
        Product p = new Product();
        p.setId(rs.getInt("id"));
        p.setName(rs.getString("name"));
        p.setDescription(rs.getString("description"));
        p.setPrice(rs.getBigDecimal("price"));
        p.setCategory(rs.getString("category"));
        p.setImageUrl(rs.getString("image_url"));
        p.setStockQuantity(rs.getInt("stock_quantity"));
        return p;
    };

    public List<Product> findAll() {
        return jdbcTemplate.query("SELECT * FROM products", rowMapper);
    }

    public Product findById(Integer id) {
        try {
            return jdbcTemplate.queryForObject("SELECT * FROM products WHERE id = ?", rowMapper, id);
        } catch (Exception e) {
            return null;
        }
    }

    public List<Product> findByCategory(String category) {
        return jdbcTemplate.query("SELECT * FROM products WHERE category = ?", rowMapper, category);
    }

    public List<Product> searchByNameOrDesc(String term) {
        String sql = "SELECT * FROM products WHERE name LIKE '%" + term + "%'";
        return jdbcTemplate.query(sql, rowMapper);
    }

    public void save(Product p) {
        jdbcTemplate.update(
                "INSERT INTO products (name, description, price, category, image_url, stock_quantity) VALUES (?, ?, ?, ?, ?, ?)",
                p.getName(), p.getDescription(), p.getPrice(), p.getCategory(), p.getImageUrl(), p.getStockQuantity());
    }

    public void deleteById(Integer id) {
        jdbcTemplate.update("DELETE FROM products WHERE id = ?", id);
    }
}
