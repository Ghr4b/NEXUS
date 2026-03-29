package com.ctf.challenge.repository;

import com.ctf.challenge.model.User;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.core.RowMapper;
import org.springframework.stereotype.Repository;

@Repository
public class UserRepository {
    @Autowired
    private JdbcTemplate jdbcTemplate;

    private final RowMapper<User> rowMapper = (rs, rowNum) -> {
        User u = new User();
        u.setId(rs.getInt("id"));
        u.setUsername(rs.getString("username"));
        u.setPassword(rs.getString("password"));
        u.setRole(rs.getString("role"));
        return u;
    };

    public User findByUsername(String username) {
        try {
            return jdbcTemplate.queryForObject("SELECT * FROM users WHERE username = ?", rowMapper, username);
        } catch (Exception e) {
            return null;
        }
    }

    public void save(User user) {
        jdbcTemplate.update("INSERT INTO users (username, password, role) VALUES (?, ?, ?)",
                user.getUsername(), user.getPassword(), user.getRole());
    }
}
