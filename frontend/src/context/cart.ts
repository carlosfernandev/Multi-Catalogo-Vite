// src/context/cart.ts
// Tipos, constantes y hook del carrito de compras separados del Provider para cumplir con Fast Refresh.
import { createContext, useContext, type ReactNode } from 'react';
import type { IProducto } from '../components/Catalogo/interfaces/iproducto';

// 1. Elemento del carrito: un producto + la cantidad seleccionada
export interface CartItem extends IProducto {
  cantidad: number;
}

export interface CartContextType {
  cart: CartItem[];
  addToCart: (producto: IProducto) => void;
  removeFromCart: (id: number) => void;
  incrementQuantity: (id: number) => void;
  decrementQuantity: (id: number) => void;
  clearCart: () => void;
  totalItems: number;
  totalPrice: number;
}

// 2. Prefijo de la clave de persistencia en localStorage.
// Se completa con el correo del usuario para que CADA CUENTA tenga su propio carrito.
export const CART_STORAGE_PREFIX = "multicatalogo_carrito_";

// 3. Creamos el contexto indicando que puede ser CartContextType o undefined
export const CartContext = createContext<CartContextType | undefined>(undefined);

// 4. Hook personalizado con validación de tipo
export const useCart = () => {
  const context = useContext(CartContext);
  if (!context) {
    throw new Error("useCart debe ser usado dentro de un CartProvider");
  }
  return context;
};

// 5. Tipamos los props del Provider
export interface CartProviderProps {
  children: ReactNode;
}
