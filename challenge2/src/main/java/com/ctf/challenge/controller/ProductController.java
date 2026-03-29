package com.ctf.challenge.controller;

import com.ctf.challenge.service.ProductService;
import com.ctf.challenge.util.SqlValidator;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestParam;

import java.util.List;
import com.ctf.challenge.model.Product;
import java.util.Comparator;

@Controller
public class ProductController {

    @Autowired
    private ProductService productService;

    @GetMapping("/catalog")
    public String browseProducts(
            @RequestParam(required = false) String category,
            @RequestParam(required = false) String search,
            @RequestParam(required = false) String minPrice,
            @RequestParam(required = false) String maxPrice,
            @RequestParam(required = false) String inStock,
            @RequestParam(required = false, defaultValue = "default") String sort,
            Model model) {

        List<Product> products;

        if (search != null && !search.isEmpty()) {
            try {
                String sanitized = SqlValidator.sanitize(search);
                products = productService.searchByNameOrDesc(sanitized);
                model.addAttribute("searchApplied", true);
                model.addAttribute("searchTerm", search);
            } catch (IllegalArgumentException e) {
                model.addAttribute("error", e.getMessage());
                products = productService.getAllProducts();
            }
        } else if (category != null && !category.isEmpty()) {
            products = productService.getProductsByCategory(category);
            model.addAttribute("selectedCategory", category);
        } else {
            products = productService.getAllProducts();
        }

        if (minPrice != null && !minPrice.isEmpty()) {
            try {
                double min = Double.parseDouble(minPrice);
                products = products.stream()
                        .filter(p -> p.getPrice().doubleValue() >= min)
                        .toList();
            } catch (NumberFormatException ignored) {
            }
        }
        if (maxPrice != null && !maxPrice.isEmpty()) {
            try {
                double max = Double.parseDouble(maxPrice);
                products = products.stream()
                        .filter(p -> p.getPrice().doubleValue() <= max)
                        .toList();
            } catch (NumberFormatException ignored) {
            }
        }
        if ("true".equals(inStock)) {
            products = products.stream()
                    .filter(p -> p.getStockQuantity() > 0)
                    .toList();
            model.addAttribute("inStockOnly", true);
        }

        products = switch (sort) {
            case "price_asc" -> products.stream().sorted(Comparator.comparing(Product::getPrice)).toList();
            case "price_desc" -> products.stream().sorted(Comparator.comparing(Product::getPrice).reversed()).toList();
            case "name_asc" -> products.stream().sorted(Comparator.comparing(p -> p.getName().toLowerCase())).toList();
            case "popular" -> products.stream().sorted(Comparator.comparing(Product::getStockQuantity)).toList();
            default -> products;
        };

        model.addAttribute("products", products);
        model.addAttribute("selectedSort", sort);
        model.addAttribute("minPrice", minPrice);
        model.addAttribute("maxPrice", maxPrice);
        return "catalog/list";
    }

    @GetMapping("/product/{id}")
    public String viewProduct(@PathVariable Integer id, Model model) {
        model.addAttribute("product", productService.getProduct(id));
        return "catalog/detail";
    }
}
