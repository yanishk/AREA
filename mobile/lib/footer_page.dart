import 'package:flutter/material.dart';

import "login_page.dart";
import "exploreAll.dart";
import "edit_account.dart";
import "create.dart";
import 'globalData.dart';

class FooterNavigationBar extends StatelessWidget {
  FooterNavigationBar({
    super.key,
    required this.currentIndex,
  });

  final int currentIndex;
  final List<Widget> _pages = [
    LoginPage(),
    ExplorePageAll(),
    CreateWorkflowsPage(),
    LoginPage(),
    EditAccountPage(),
  ];

  @override
  Widget build(BuildContext context) {
    return BottomNavigationBar (
      type: BottomNavigationBarType.fixed,
      showUnselectedLabels: true,
      selectedItemColor: Colors.black,
      unselectedItemColor: Colors.black45,
      backgroundColor: Colors.white,
      selectedLabelStyle: const TextStyle(fontSize: 15, fontWeight: FontWeight.bold),
      unselectedLabelStyle: const TextStyle(fontSize: 12, fontWeight: FontWeight.bold),
      currentIndex: currentIndex,
      onTap: (int index) {
        if (index != currentIndex) {
          Navigator.push(
            Globaldata.myContext,
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