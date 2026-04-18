package com.ctf.challenge.service;

import com.ctf.challenge.model.Product;
import com.ctf.challenge.repository.ProductRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class ProductService {

    @Autowired
    private ProductRepository productRepository;

    public List<Product> getAllProducts() {
        return productRepository.findAll();
    }

    public List<Product> getProductsByCategory(String category) {
        if (category == null || category.isEmpty()) {
            return getAllProducts();
        }
        return productRepository.findByCategory(category);
    }

    public Product getProduct(Integer id) {
        return productRepository.findById(id);
    }

    public List<Product> searchByNameOrDesc(String term) {
        return productRepository.searchByNameOrDesc(term);
    }
}
