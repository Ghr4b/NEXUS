package com.ctf.challenge.service;

import com.ctf.challenge.model.Order;
import com.ctf.challenge.model.OrderItem;
import com.ctf.challenge.model.Product;
import com.ctf.challenge.repository.OrderRepository;
import com.ctf.challenge.repository.ProductRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;
import java.util.List;
import java.util.Map;

@Service
public class OrderService {

    @Autowired
    private OrderRepository orderRepository;

    @Autowired
    private ProductRepository productRepository;

    @Transactional
    public Integer checkout(Integer userId, Map<Integer, Integer> cart, BigDecimal discountFactor) {
        if (cart == null || cart.isEmpty())
            return null;

        BigDecimal total = BigDecimal.ZERO;
        for (Map.Entry<Integer, Integer> entry : cart.entrySet()) {
            Product p = productRepository.findById(entry.getKey());
            if (p != null) {
                total = total.add(p.getPrice().multiply(new BigDecimal(entry.getValue())));
            }
        }

        total = total.multiply(discountFactor);

        Order order = new Order();
        order.setUserId(userId);
        order.setStatus("PENDING");
        order.setTotalPrice(total);
        Integer orderId = orderRepository.save(order);

        for (Map.Entry<Integer, Integer> entry : cart.entrySet()) {
            Product p = productRepository.findById(entry.getKey());
            if (p != null) {
                OrderItem item = new OrderItem();
                item.setOrderId(orderId);
                item.setProductId(p.getId());
                item.setQuantity(entry.getValue());
                item.setPrice(p.getPrice());
                orderRepository.saveItem(item);

                p.setStockQuantity(p.getStockQuantity() - entry.getValue());
                productRepository.save(p);
            }
        }

        return orderId;
    }

    public List<Order> getUserOrders(Integer userId) {
        return orderRepository.findByUserId(userId);
    }

    public Order getOrderDetails(Integer orderId) {
        return orderRepository.findById(orderId);
    }

    public List<OrderItem> getOrderItems(Integer orderId) {
        return orderRepository.findItemsByOrderId(orderId);
    }

    public List<Order> getAllOrders() {
        return orderRepository.findAll();
    }

    public void updateStatus(Integer orderId, String status) {
        orderRepository.updateStatus(orderId, status);
    }
}
