import 'package:flutter/material.dart';

import "login_page.dart";
import "sign_up.dart";
import "explore.dart";
import "edit_account.dart";

class FooterNavigationBar extends StatelessWidget {
  FooterNavigationBar({
    super.key,
    required this.currentIndex,
  });

  final int currentIndex;
  final List<Widget> _pages = [
    LoginPage(),
    ExplorePage(),
    SignUp(),
    LoginPage(),
    EditAccount(),
  ];

  @override
  Widget build(BuildContext context) {
    return BottomNavigationBar (
      type: BottomNavigationBarType.fixed,
      showUnselectedLabels: true,
      selectedItemColor: Colors.black,
      unselectedItemColor: Colors.black45,
      backgroundColor: Colors.white,
      selectedLabelStyle: TextStyle(fontSize: 15, fontWeight: FontWeight.bold),
      unselectedLabelStyle: TextStyle(fontSize: 12, fontWeight: FontWeight.bold),
      currentIndex: currentIndex,
      onTap: (int index) {
        if (index != currentIndex) {
          Navigator.push(
            context,
            MaterialPageRoute(builder: (context) => _pages[index]),
          );
        }
      },
      items: const <BottomNavigationBarItem>[
        BottomNavigationBarItem(
          icon: Icon(Icons.call_to_action_outlined),
          label: 'My Workflow',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.search),
          label: 'Explore',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.add_circle_outline_outlined),
          label: 'Create',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.format_list_bulleted_outlined),
          label: 'Activity',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.person_2_outlined),
          label: 'Profile',
        ),
      ],
    );
  }
}