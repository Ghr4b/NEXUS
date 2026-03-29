package com.ctf.challenge.repository;

import com.ctf.challenge.model.PromoCode;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.core.RowMapper;
import org.springframework.stereotype.Repository;

import java.sql.Types;
import java.util.List;

@Repository
public class PromoRepository {

    @Autowired
    private JdbcTemplate jdbcTemplate;

    private final RowMapper<PromoCode> rowMapper = (rs, rowNum) -> {
        PromoCode p = new PromoCode();
        p.setId(rs.getInt("id"));
        p.setCode(rs.getString("code"));
        p.setDiscountPercentage(rs.getInt("discount_percentage"));
        p.setIsActive(rs.getBoolean("is_active"));
        if (rs.getObject("user_id") != null) {
            p.setUserId(rs.getInt("user_id"));
        }
        return p;
    };

    public PromoCode findByCode(String code) {
        try {
            return jdbcTemplate.queryForObject("SELECT * FROM promo_codes WHERE code = ?", rowMapper, code);
        } catch (Exception e) {
            return null;
        }
    }

    public List<PromoCode> findAll() {
        return jdbcTemplate.query("SELECT * FROM promo_codes", rowMapper);
    }

    public void save(PromoCode promo) {
        if (promo.getUserId() == null) {
            jdbcTemplate.update(
                    "INSERT INTO promo_codes (code, discount_percentage, is_active, user_id) VALUES (?, ?, ?, NULL)",
                    promo.getCode(), promo.getDiscountPercentage(), promo.getIsActive());
        } else {
            jdbcTemplate.update(
                    "INSERT INTO promo_codes (code, discount_percentage, is_active, user_id) VALUES (?, ?, ?, ?)",
                    promo.getCode(), promo.getDiscountPercentage(), promo.getIsActive(), promo.getUserId());
        }
    }

    public void deactivate(Integer id) {
        jdbcTemplate.update("UPDATE promo_codes SET is_active = FALSE WHERE id = ?", id);
    }
}
