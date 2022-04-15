insert into apis (method, route, name, status) values ("GET","/organizations","获取组织列表",1),("GET","/organizations/:id","根据ID获取组织",1),("PUT","/organizations/:id","更新组织",1),("POST","/organizations","新建组织",1),("GET","/projects","获取项目列表",1),("GET","/projects/:id","根据ID获取项目",1),("PUT","/projects/:id","更新项目",1),("POST","/projects","新建项目",1),("DELETE","/projects/:id","删除项目",1),("GET","/events","获取事件列表",1),("GET","/events/:id","根据ID获取事件",1),("PUT","/events/:id","更新事件",1),("GET","/components","获取组件列表",1),("GET","/components/:id","根据ID获取组件",1),("GET","/roles","获取角色列表",1),("GET","/roles/:id","根据ID获取角色",1),("PUT","/roles/:id","更新角色",1),("POST","/roles","新建角色",1),("DELETE","/roles/:id","删除角色",1),("PUT","/users/:id","更新用户信息",1),("GET","/users","获取用户列表",1),("GET","/users/:id","根据ID获取用户",1),("GET","/apis","获取API列表",1),("GET","/apis/:id","根据ID获取API",1),("PUT","/apis/:id","更新API",1),("POST","/apis","新建API",1),("GET","/menus","获取菜单列表",1),("GET","/menus/:id","根据ID获取菜单",1),("POST","/menus","新建菜单",1),("PUT","/menus/:id","更新菜单",1),("DELETE","/menus/:id","删除菜单",1),("GET","/rolemenus/:id","获取角色菜单绑定关系",1),("POST","/rolemenus/:id","更新角色菜单绑定关系",1),("GET","/menuapis/:id","获取菜单API绑定关系",1),("POST","/menuapis/:id","更新菜单API绑定关系",1),("GET","/mymenu","获取当前用户的菜单",1),("GET","/clients","获取客户列表",1),("GET","/clients/:id","根据ID获取客户",1),("PUT","/clients/:id","更新客户",1),("POST","/clients","新建客户",1),("GET","/positions","获取职位列表",1),("GET","/positions/:id","根据ID获取职位",1),("PUT","/positions/:id","更新职位",1),("POST","/positions","新建职位",1),("GET","/members","获取项目成员",1),("POST","/members","更新项目成员",1),("GET","/templates","获取模板列表",1),("GET","/templates/:id","根据ID获取模板",1),("PUT","/templates/:id","更新模板",1),("POST","/templates","新建模板",1),("DELETE","/templates/:id","删除模板",1),("GET","/nodes","获取节点列表",1),("GET","/nodes/:id","根据ID获取节点",1),("PUT","/nodes/:id","更新节点",1),("POST","/nodes","新建节点",1),("DELETE","/nodes/:id","删除节点",1),("GET","/elements","获取元素列表",1),("GET","/elements/:id","根据ID获取元素",1),("PUT","/elements/:id","更新元素",1),("POST","/elements","新建元素",1),("DELETE","/elements/:id","删除元素",1),("GET","/uploads","获取上传列表",1),("POST","/uploads","上传文件",1);

insert into role_menus (role_id, menu_id, status) values (1,1,1),(1,2,1),(1,3,1),(1,4,1),(1,5,1),(1,6,1),(1,7,1),(1,8,1),(1,9,1),(1,10,1),(1,11,1);

